/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
	"flag"
	"fmt"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//lo := meta_v1.ListOptions{LabelSelector: "pg-cluster=dinner"}
	//cmap, err := clientset.CoreV1().ConfigMaps("demo").Get("pgo-pgbackrest-config", meta_v1.GetOptions{})
	cmap, err := clientset.CoreV1().ConfigMaps("demo").Get("pgo-config", meta_v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	if cmap.Data["pgbackrest.conf"] != "" {
		fmt.Printf("found pgbackrest.conf in map %s", cmap.Data["pgbackrest.conf"])
	} else {
		fmt.Println("NOT found pgbackrest.conf in map")
	}

	//for _, m := range cmaps.Items {
	//}
}
