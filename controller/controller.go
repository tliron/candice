package controller

import (
	contextpkg "context"
	"fmt"
	"time"

	candiceclientset "github.com/tliron/candice/apis/clientset/versioned"
	candiceinformers "github.com/tliron/candice/apis/informers/externalversions"
	candicelisters "github.com/tliron/candice/apis/listers/candice.puccini.cloud/v1alpha1"
	"github.com/tliron/candice/client"
	candiceresources "github.com/tliron/candice/resources/candice.puccini.cloud/v1alpha1"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/logging"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	dynamicpkg "k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

//
// Controller
//

type Controller struct {
	Config      *restpkg.Config
	Dynamic     *kubernetesutil.Dynamic
	Kubernetes  kubernetes.Interface
	Candice     candiceclientset.Interface
	Client      *client.Client
	CachePath   string
	StopChannel <-chan struct{}

	Processors *kubernetesutil.Processors
	Events     record.EventRecorder

	KubernetesInformerFactory informers.SharedInformerFactory
	CandiceInformerFactory    candiceinformers.SharedInformerFactory

	Devices candicelisters.DeviceLister

	Context contextpkg.Context
	Log     logging.Logger
}

func NewController(context contextpkg.Context, toolName string, clusterMode bool, clusterRole string, namespace string, dynamic dynamicpkg.Interface, kubernetes kubernetes.Interface, apiExtensions apiextensionspkg.Interface, candice candiceclientset.Interface, config *restpkg.Config, informerResyncPeriod time.Duration, stopChannel <-chan struct{}) *Controller {
	if clusterMode {
		namespace = ""
		if clusterRole != "" {
			clusterRole = "cluster-admin"
		}
	}

	log := logging.GetLoggerf("%s.controller", toolName)

	self := Controller{
		Config:      config,
		Dynamic:     kubernetesutil.NewDynamic(toolName, dynamic, kubernetes.Discovery(), namespace, context),
		Kubernetes:  kubernetes,
		Candice:     candice,
		StopChannel: stopChannel,
		Processors:  kubernetesutil.NewProcessors(toolName),
		Events:      kubernetesutil.CreateEventRecorder(kubernetes, "Candice", log),
		Context:     context,
		Log:         log,
	}

	self.Client = client.NewClient(
		kubernetes,
		apiExtensions,
		candice,
		kubernetes.CoreV1().RESTClient(),
		config,
		context,
		clusterMode,
		clusterRole,
		namespace,
		NamePrefix,
		PartOf,
		ManagedBy,
		OperatorImageReference,
		fmt.Sprintf("%s.client", toolName),
	)

	if clusterMode {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactory(kubernetes, informerResyncPeriod)
		self.CandiceInformerFactory = candiceinformers.NewSharedInformerFactory(candice, informerResyncPeriod)
	} else {
		self.KubernetesInformerFactory = informers.NewSharedInformerFactoryWithOptions(kubernetes, informerResyncPeriod, informers.WithNamespace(namespace))
		self.CandiceInformerFactory = candiceinformers.NewSharedInformerFactoryWithOptions(candice, informerResyncPeriod, candiceinformers.WithNamespace(namespace))
	}

	// Informers
	deviceInformer := self.CandiceInformerFactory.Candice().V1alpha1().Devices()

	// Listers
	self.Devices = deviceInformer.Lister()

	// Processors

	processorPeriod := 5 * time.Second

	self.Processors.Add(candiceresources.DeviceGVK, kubernetesutil.NewProcessor(
		toolName,
		"devices",
		deviceInformer.Informer(),
		processorPeriod,
		func(name string, namespace string) (interface{}, error) {
			return self.Client.GetDevice(namespace, name)
		},
		func(object interface{}) (bool, error) {
			return self.processDevice(object.(*candiceresources.Device))
		},
	))

	return &self
}

func (self *Controller) Run(concurrency uint, startup func()) error {
	defer utilruntime.HandleCrash()

	self.Log.Info("starting informer factories")
	self.KubernetesInformerFactory.Start(self.StopChannel)
	self.CandiceInformerFactory.Start(self.StopChannel)

	self.Log.Info("waiting for processor informer caches to sync")
	utilruntime.HandleError(self.Processors.WaitForCacheSync(self.StopChannel))

	self.Log.Infof("starting processors (concurrency=%d)", concurrency)
	self.Processors.Start(concurrency, self.StopChannel)
	defer self.Processors.ShutDown()

	if startup != nil {
		go startup()
	}

	<-self.StopChannel

	self.Log.Info("shutting down")

	return nil
}
