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

// Package main for a test operator
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	testoperator "github.com/imcsk8/k8s-operators/pkg/apis/testoperator/v1alpha1"
	testoperatorclient "github.com/imcsk8/k8s-operators/pkg/client/clientset/versioned/typed/testoperator/v1alpha1"
	opkit "github.com/imcsk8/operator-kit"
	"k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Getting kubernetes context")
	context, testOperatorClientset, err := createContext()
	if err != nil {
		fmt.Printf("failed to create context. %+v\n", err)
		os.Exit(1)
	}

	// Create and wait for CRD resources
	fmt.Println("Registering the sample resource")
	resources := []opkit.CustomResource{testoperator.TestOperatorResource}
	err = opkit.CreateCustomResources(*context, resources)
	if err != nil {
		fmt.Printf("failed to create custom resource. %+v\n", err)
		os.Exit(1)
	}

	// create signals to stop watching the resources
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start watching the sample resource
	fmt.Println("Watching the test operator resource")
	controller := newTestOperatorController(context, sampleClientset)
	controller.StartWatch(v1.NamespaceAll, stopChan)

	for {
		select {
		case <-signalChan:
			fmt.Println("shutdown signal received, exiting...")
			close(stopChan)
			return
		}
	}
}

func createContext() (*opkit.Context, testoperatorclient.MyprojectV1alpha1Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
	}

	testOperatorClientset, err := testoperatorclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create sample clientset. %+v", err)
	}

	context := &opkit.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		Interval:              500 * time.Millisecond,
		Timeout:               60 * time.Second,
	}
	return context, testOperatorClientset, nil

}
