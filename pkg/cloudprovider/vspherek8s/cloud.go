/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vspherek8s

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/cloudprovider/vspherek8s/loadbalancer"
	"github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/nsxt"
	"k8s.io/client-go/rest"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"

	lcfg "github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/cloudprovider/vspherek8s/loadbalancer/config"
	ncfg "github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/nsxt/config"
)

const (
	// RegisteredProviderName is the name of the cloud provider registered with
	// Kubernetes.
	RegisteredProviderName string = "vsphere-with-k8s"

	// ProviderName is the name used for constructing Provider ID
	ProviderName string = "vsphere"

	clientName string = "vsphere-k8s-cloud-controller-manager"
)

func init() {
	cloudprovider.RegisterCloudProvider(RegisteredProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		if config == nil {
			return nil, fmt.Errorf("no vsphere-with-k8s cloud provider config file given")
		}

		byConfig, err := ioutil.ReadAll(config)
		if err != nil {
			klog.Errorf("ReadAll failed: %s", err)
			return nil, err
		}

		cfg, err := readConfig(byConfig)
		if err != nil {
			// we got an error where the decode wasn't related to a missing type
			return nil, err
		}

		nsxtcfg, err := ncfg.ReadNsxtConfig(byConfig)
		if err != nil {
			klog.Errorf("ReadNsxtConfig failed: %s", err)
			nsxtcfg = nil
		}
		lbcfg, err := lcfg.ReadLBConfig(byConfig)
		if err != nil {
			klog.Errorf("ReadLBConfig failed: %s", err)
			lbcfg = nil
		}
		return newVSphereWithK8s(cfg, nsxtcfg, lbcfg)
	})
}

// Creates new Controller node interface and returns
func newVSphereWithK8s(cfg *config, nsxtcfg *ncfg.NsxtConfig, lbcfg *lcfg.LBConfig) (*VSphereWithK8s, error) {
	cp := &VSphereWithK8s{
		config:     *cfg,
		nsxtconfig: nsxtcfg,
		lbconfig:   lbcfg,
	}

	if cp.isLoadBalancerSupportEnabled() {
		ncm, err := nsxt.NewConnectorManager(nsxtcfg)
		if err != nil {
			return nil, err
		}

		lb, err := loadbalancer.NewLBProvider(lbcfg, ncm.GetConnector())
		if err != nil {
			return nil, err
		}

		cp.nsxtConnectorMgr = ncm
		cp.loadbalancer = lb
	}

	return cp, nil
}

func (cp *VSphereWithK8s) isLoadBalancerSupportEnabled() bool {
	return cp.lbconfig != nil && cp.nsxtconfig != nil && cp.lbconfig.LoadBalancer.Enabled
}

// Initialize initializes the vSphere with k8s cloud provider.
func (cp *VSphereWithK8s) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	klog.V(0).Info("Initializing vSphere with Kubernetes Cloud Provider")

	client, err := clientBuilder.Client(clientName)
	if err != nil {
		klog.Fatalf("Failed to create cloud provider client: %v", err)
	}

	cp.client = client

	kcfg, err := cp.getRestConfig()
	if err != nil {
		klog.Fatalf("Failed to create rest config to communicate with supervisor: %v", err)
	}

	instances, err := NewInstancesV2(cp.config.Supervisor.Namespace, kcfg)
	if err != nil {
		klog.Errorf("Failed to init Instance: %v", err)
	}
	cp.instances = instances

	if cp.isLoadBalancerSupportEnabled() {
		klog.Info("initializing load balancer support")
		if loadbalancer.ClusterName == "" {
			klog.Warning("Missing cluster id, no periodical cleanup possible")
		}
		cp.loadbalancer.Initialize(loadbalancer.ClusterName, client, stop)
	}

	klog.V(0).Info("Initializing vSphere with Kubernetes Cloud Provider Succeeded")
}

// LoadBalancer returns a balancer interface. Also returns true if the
// interface is supported, false otherwise.
func (cp *VSphereWithK8s) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	if cp.isLoadBalancerSupportEnabled() {
		return cp.loadbalancer, true
	}
	klog.V(1).Info("The vSphere with Kubernetes cloud provider has been configured without load balances")
	return nil, false
}

// Instances returns an instances interface. Also returns true if the
// interface is supported, false otherwise.
func (cp *VSphereWithK8s) Instances() (cloudprovider.Instances, bool) {
	klog.V(1).Info("The vSphere with Kubernetes cloud provider does only support instancesV2 not instances")
	return nil, false
}

// InstancesV2 returns an implementation of cloudprovider.InstancesV2.
func (cp *VSphereWithK8s) InstancesV2() (cloudprovider.InstancesV2, bool) {
	klog.V(1).Info("Enabling Instances interface on vSphere with Kubernetes cloud provider")
	return cp.instances, true
}

// Zones returns a zones interface. Also returns true if the interface
// is supported, false otherwise.
func (cp *VSphereWithK8s) Zones() (cloudprovider.Zones, bool) {
	klog.V(1).Info("The vSphere with Kubernetes cloud provider does not support zones")
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface
// is supported, false otherwise.
func (cp *VSphereWithK8s) Clusters() (cloudprovider.Clusters, bool) {
	klog.V(1).Info("The vSphere with Kubernetes cloud provider does not support clusters")
	return nil, false
}

// Routes returns a routes interface along with whether the interface
// is supported.
func (cp *VSphereWithK8s) Routes() (cloudprovider.Routes, bool) {
	klog.V(1).Info("The vSphere with Kubernetes cloud provider does not support routes")
	return nil, false
}

// ProviderName returns the cloud provider ID.
// Note: Returns 'vsphere' instead of 'vsphere-with-k8s'
// since CAPV expects the ProviderID to be in form 'vsphere://***'
// https://github.com/kubernetes/cloud-provider-vsphere/issues/447
func (cp *VSphereWithK8s) ProviderName() string {
	return ProviderName
}

// HasClusterID returns true if a ClusterID is required and set/
func (cp *VSphereWithK8s) HasClusterID() bool {
	return true
}

func (cp *VSphereWithK8s) getRestConfig() (*rest.Config, error) {
	var caData []byte
	if cp.config.Supervisor.CAData != "" {
		bytes, err := base64.StdEncoding.DecodeString(cp.config.Supervisor.CAData)
		if err != nil {
			return nil, fmt.Errorf("cannot decode config.Supervisor.CAData: %w", err)
		}
		caData = bytes
	}
	return &rest.Config{
		Host: cp.config.Supervisor.Apiserver,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:     caData,
			ServerName: cp.config.Supervisor.ApiserverFQDN,
			Insecure:   cp.config.Supervisor.Insecure,
		},
		BearerToken: string(cp.config.Supervisor.Token),
	}, nil
}
