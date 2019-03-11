package main

import (
	"flag"
	"fmt"
	"github.com/crunchydata/postgres-operator/kubeapi"
	"k8s.io/api/storage/v1"
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

	var scList *v1.StorageClassList
	scList, err = kubeapi.GetStorageClasses(clientset)
	if err != nil {
		return
	}

	fmt.Printf("scList len %d\n", len(scList.Items))

	for _, j := range scList.Items {
		fmt.Printf("storage class %s\n", j.Name)
	}

}
