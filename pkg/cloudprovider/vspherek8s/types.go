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
	"github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/cloudprovider/vspherek8s/loadbalancer"
	lcfg "github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/cloudprovider/vspherek8s/loadbalancer/config"
	"github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/nsxt"
	ncfg "github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/nsxt/config"
	clientset "k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
)

// VSphereWithK8s is an implementation of cloud provider Interface for vsphere with kubernetes.
type VSphereWithK8s struct {
	config
	client    clientset.Interface
	instances cloudprovider.InstancesV2

	nsxtconfig       *ncfg.NsxtConfig
	lbconfig         *lcfg.LBConfig
	loadbalancer     loadbalancer.LBProvider
	nsxtConnectorMgr *nsxt.ConnectorManager
}
