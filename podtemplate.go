package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"os"
	"text/template"
)

const PodTemplateFile = "pod-template.yaml"

type PodFields struct {
	PodName string
}

func main() {
	buf, err := ioutil.ReadFile(PodTemplateFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	value := string(buf)
	fmt.Println(value)
	var tmpl *template.Template
	tmpl = template.Must(template.New("podtemplate").Parse(value))
	if tmpl == nil {
		fmt.Println("error in template")
	}

	myPodInfo := PodFields{
		PodName: "rqpod1",
	}
	tmpl.Execute(os.Stdout, myPodInfo)

	var podBuffer bytes.Buffer
	tmpl.Execute(&podBuffer, myPodInfo)

	fmt.Printf("final string is %s\n", podBuffer.String())

	mypod := corev1.Pod{}
	err = yaml.Unmarshal(podBuffer.Bytes(), &mypod)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	fmt.Println("completed ok")
}
