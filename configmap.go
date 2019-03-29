package main

import (
	"flag"
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
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
	namespace := "pgo"
	configmapname := "pgo-config"
	key := "pool_hba.conf"
	cmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(configmapname, meta_v1.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	fmt.Println("yes pgo-config configmap is there")

	if cmap.Data[key] != "" {
		fmt.Printf("found %s key in map %s", key, configmapname)
		fmt.Printf("value is %s", cmap.Data[key])
	} else {
		fmt.Println("%s key NOT found in map %s", key, configmapname)
	}

	//for _, m := range cmaps.Items {
	//}
}
