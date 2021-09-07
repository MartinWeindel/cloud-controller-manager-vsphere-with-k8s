module github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s

go 1.16

require (
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/vmware-tanzu/vm-operator-api v0.1.4-0.20201118171008-5ca641b0e126
	github.com/vmware/vsphere-automation-sdk-go/lib v0.2.0
	github.com/vmware/vsphere-automation-sdk-go/runtime v0.2.0
	github.com/vmware/vsphere-automation-sdk-go/services/nsxt v0.3.0
	gopkg.in/gcfg.v1 v1.2.3
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/cloud-provider v0.21.2
	k8s.io/component-base v0.21.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.6.5
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
