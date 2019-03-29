package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/crunchydata/postgres-operator/kubeapi"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()

	log.SetLevel(log.DebugLevel)

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

	jobName := "backup-donk-kiyn"
	namespace := "pgouser1"

	job, found := kubeapi.GetJob(clientset, jobName, namespace)
	if !found {
		log.Error("not found")
		return
	}

	if job.Status.Succeeded > 0 {
		log.Debug("job succeeded")

	}

	updateBackupPaths(clientset, jobName, namespace)

}

func updateBackupPaths(clientset *kubernetes.Clientset, jobName, namespace string) {
	//its pod has this label job-name=backup-yank-fjus
	selector := "job-name=" + jobName
	log.Debugf("looking for pod with selector %s", selector)
	podList, err := kubeapi.GetPods(clientset, selector, namespace)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if len(podList.Items) != 1 {
		log.Error("could not find a pod for this job")
		return
	}

	podName := podList.Items[0].Name
	log.Debugf("found pod %s", podName)
	backupPath, err := getBackupPath(clientset, podName, namespace)
	if err != nil {
		log.Error("error in getting logs %s", err.Error())
		return
	}
	log.Debugf("backupPath is %s", backupPath)

}

//this func assumes the pod has completed and its a backup job pod
func getBackupPath(clientset *kubernetes.Clientset, podName, namespace string) (string, error) {
	opts := v1.PodLogOptions{
		Container: "backup",
	}
	var logs bytes.Buffer

	err := kubeapi.GetLogs(clientset, opts, &logs, podName, namespace)
	if err != nil {
		log.Error("error in getting logs %s", err.Error())
		return "", err
	}

	//this is what the backup container puts in its log and what
	//we are looking to parse out of its container log
	token := "BACKUP_PATH is set to /pgdata/"

	var backupPath string
	scanner := bufio.NewScanner(strings.NewReader(logs.String()))
	for scanner.Scan() {
		rawStr := scanner.Text()
		rawlen := len(rawStr)
		idx := strings.Index(rawStr, token)
		if idx > -1 {
			log.Debugf("log line of interest %s", scanner.Text())
			content := rawlen - idx
			log.Debugf("raw length %d token at %d content at %d\n", rawlen, idx, content)
			//parsed := rawStr[idx+len(token):]
			parsed := rawStr[content:]
			backupPath = parsed[:rawlen-content-5]
		}
	}

	return backupPath, nil

}
