package main

import (
	"flag"
	"fmt"
	"github.com/crunchydata/postgres-operator/kubeapi"
	"github.com/crunchydata/postgres-operator/util"
	apiv1 "k8s.io/api/batch/v1"
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

	clusterName := "fish"
	namespace := "demo"

	var jobList *apiv1.JobList
	jobList, err = kubeapi.GetJobs(clientset, util.LABEL_PG_CLUSTER+"="+clusterName, namespace)
	if err != nil {
		return
	}

	for _, j := range jobList.Items {
		if j.Status.Succeeded > 0 {
			fmt.Printf("removing Job %s since it was completed \n", j.Name)
			err := kubeapi.DeleteJob(clientset, j.Name, namespace)
			if err != nil {
				fmt.Println(err)
			}

		}
	}

}
