module github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/vm-operator-api v0.1.4-0.20201118171008-5ca641b0e126
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/cloud-provider v0.21.2
	k8s.io/component-base v0.21.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.6.5
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
