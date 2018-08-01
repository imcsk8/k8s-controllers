/*
Copyright 2016 The Rook Authors. All rights reserved.
Copyright 2018 Iván Chavero. All rights reserved.

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

// Package main for a sample operator
package main

import (
	"fmt"

	testoperator "github.com/imcsk8/k8s-operators/pkg/apis/testoperator/v1alpha1"
	testoperatorclient "github.com/imcsk8/k8s-operators/pkg/client/clientset/versioned/typed/testoperator/v1alpha1"
	opkit "github.com/rook/operator-kit"
	"k8s.io/client-go/tools/cache"
)

// TestOperatorController represents a controller object for sample custom resources
type TestOperatorController struct {
	context               *opkit.Context
	testOperatorClientset testoperatorclient.MyprojectV1alpha1Interface
}

// newTestOperatorController create controller for watching sample custom resources created
func newTestOperatorController(context *opkit.Context, testOperatorClientset testoperatorclient.MyprojectV1alpha1Interface) *TestOperatorController {
	return &TestOperatorController{
		context:               context,
		testOperatorClientset: testOperatorClientset,
	}
}

// Watch watches for instances of TestOperator custom resources and acts on them
func (c *TestOperatorController) StartWatch(namespace string, stopCh chan struct{}) error {

	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.testOperatorClientset.RESTClient()
	watcher := opkit.NewWatcher(testoperator.TestOperatorResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&sample.Sample{}, stopCh)
	return nil
}

func (c *TestOperatorController) onAdd(obj interface{}) {
	s := obj.(*testoperator.TestOperator).DeepCopy()

	fmt.Printf("Added TestOperator '%s' with Hello=%s\n", s.Name, s.Spec.Hello)
}

func (c *TestOperatorController) onUpdate(oldObj, newObj interface{}) {
	oldTestOperator := oldObj.(*testoperator.TestOperator).DeepCopy()
	newTestOperator := newObj.(*testoperator.TestOperator).DeepCopy()

	fmt.Printf("Updated TestOperator '%s' from %s to %s\n", newTestOperator.Name, oldTestOperator.Spec.Hello, newTestOperator.Spec.Hello)
}

func (c *TestOperatorController) onDelete(obj interface{}) {
	s := obj.(*testoperator.TestOperator).DeepCopy()

	fmt.Printf("Deleted testoperator '%s' with Hello=%s\n", s.Name, s.Spec.Hello)
}
