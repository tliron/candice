package commands

import (
	contextpkg "context"
	"fmt"

	candicepkg "github.com/tliron/candice/apis/clientset/versioned"
	"github.com/tliron/candice/client"
	"github.com/tliron/candice/controller"
	kubernetesutil "github.com/tliron/kutil/kubernetes"
	"github.com/tliron/kutil/util"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	restpkg "k8s.io/client-go/rest"
)

//
// Client
//

type Client struct {
	Config     *restpkg.Config
	Kubernetes kubernetespkg.Interface
	REST       restpkg.Interface
	Context    contextpkg.Context
	Namespace  string
}

func NewClient() *Client {
	config, err := kubernetesutil.NewConfigFromFlags(masterUrl, kubeconfigPath, kubeconfigContext, log)
	util.FailOnError(err)

	kubernetes, err := kubernetespkg.NewForConfig(config)
	util.FailOnError(err)

	namespace_ := namespace
	if clusterMode {
		namespace_ = ""
	} else if namespace_ == "" {
		if namespace__, ok := kubernetesutil.GetConfiguredNamespace(kubeconfigPath, kubeconfigContext); ok {
			namespace_ = namespace__
		}
		if namespace_ == "" {
			util.Fail("could not discover namespace and \"--namespace\" not provided")
		}
	}

	return &Client{
		Config:     config,
		Kubernetes: kubernetes,
		REST:       kubernetes.CoreV1().RESTClient(),
		Context:    context,
		Namespace:  namespace_,
	}
}

func (self *Client) Client() *client.Client {
	apiExtensions, err := apiextensionspkg.NewForConfig(self.Config)
	util.FailOnError(err)

	candice, err := candicepkg.NewForConfig(self.Config)
	util.FailOnError(err)

	return client.NewClient(
		self.Kubernetes,
		apiExtensions,
		candice,
		self.REST,
		self.Config,
		self.Context,
		clusterMode,
		clusterRole,
		self.Namespace,
		controller.NamePrefix,
		controller.PartOf,
		controller.ManagedBy,
		controller.OperatorImageReference,
		fmt.Sprintf("%s.client", toolName),
	)
}
