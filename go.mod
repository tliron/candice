module github.com/tliron/candice

go 1.15

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

require (
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/jetstack/cert-manager v1.2.0
	github.com/openconfig/goyang v0.2.4
	github.com/spf13/cobra v1.1.3
	github.com/tebeka/atexit v0.3.0
	github.com/tliron/kutil v0.1.20
	github.com/tliron/yamlkeys v1.3.5
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)
