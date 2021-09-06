/*
 * Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package vspherek8s

import (
	"context"

	vmopv1alpha1 "github.com/vmware-tanzu/vm-operator-api/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// AnnotationZone is the annotation key for the region
	AnnotationZone = "vsphere.provider.gardener.cloud/region"
	// AnnotationRegion is the annotation key for the zone
	AnnotationRegion = "vsphere.provider.gardener.cloud/zone"
)

type instances struct {
	vmClient  client.Client
	namespace string
}

// GetVmopClient gets a vm-operator-api client
// This is separate from NewVMService so that a fake client can be injected for testing
func GetVmopClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = vmopv1alpha1.AddToScheme(scheme)
	client, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	return client, err
}

type instancesv2 struct {
	vmopClient          client.Client
	supervisorNamespace string
}

func NewInstancesV2(supervisorNamespace string, kcfg *rest.Config) (cloudprovider.InstancesV2, error) {
	client, err := GetVmopClient(kcfg)
	if err != nil {
		return nil, err
	}

	return &instancesv2{vmopClient: client, supervisorNamespace: supervisorNamespace}, nil
}

// InstanceExists returns true if the instance for the given node exists according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i *instancesv2) InstanceExists(ctx context.Context, node *corev1.Node) (bool, error) {
	klog.V(4).Infof("InstanceExists for %s", node.Name)
	_, err := i.getVM(ctx, node.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	klog.V(4).Infof("InstanceExists for %s: true", node.Name)
	return true, nil
}

func (i *instancesv2) getVM(ctx context.Context, nodeName string) (*vmopv1alpha1.VirtualMachine, error) {
	key := client.ObjectKey{
		Namespace: i.supervisorNamespace,
		Name:      nodeName,
	}
	obj := &vmopv1alpha1.VirtualMachine{}
	err := i.vmopClient.Get(ctx, key, obj)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.V(4).Infof("getVM %s not found", nodeName)
		} else {
			klog.V(4).ErrorS(err, "getVM failed")
		}
		return nil, err
	}
	return obj, nil
}

// InstanceShutdown returns true if the instance is shutdown according to the cloud provider.
// Use the node.name or node.spec.providerID field to find the node in the cloud provider.
func (i *instancesv2) InstanceShutdown(ctx context.Context, node *corev1.Node) (bool, error) {
	klog.V(4).Infof("InstanceShutdown for %s", node.Name)
	vm, err := i.getVM(ctx, node.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}
	shutdown := vm.Status.PowerState != vmopv1alpha1.VirtualMachinePoweredOn
	klog.V(4).Infof("InstanceShutdown for %s: %t", node.Name, shutdown)
	return shutdown, nil
}

// InstanceMetadata returns the instance's metadata. The values returned in InstanceMetadata are
// translated into specific fields and labels in the Node object on registration.
// Implementations should always check node.spec.providerID first when trying to discover the instance
// for a given node. In cases where node.spec.providerID is empty, implementations can use other
// properties of the node like its name, labels and annotations.
func (i *instancesv2) InstanceMetadata(ctx context.Context, node *corev1.Node) (*cloudprovider.InstanceMetadata, error) {
	klog.V(4).Infof("InstanceMetadata for %s", node.Name)
	vm, err := i.getVM(ctx, node.Name)
	if err != nil {
		return nil, err
	}

	var region, zone string
	if annotations := vm.GetAnnotations(); annotations != nil {
		region = annotations[AnnotationRegion]
		zone = annotations[AnnotationZone]
	}

	metadata := &cloudprovider.InstanceMetadata{
		ProviderID:   node.Spec.ProviderID,
		InstanceType: vm.Spec.ClassName,
		//NodeAddresses: nil,
		Zone:   zone,
		Region: region,
	}
	klog.V(4).Infof("InstanceMetadata for %s: %v", node.Name, metadata)
	return metadata, nil
}
