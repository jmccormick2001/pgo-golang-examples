package main

import (
	"bytes"
	"encoding/json"
	"github.com/crunchydata/postgres-operator/apiserver"
	"github.com/crunchydata/postgres-operator/config"
	"github.com/crunchydata/postgres-operator/kubeapi"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

const PGO_ROLE = "pgo-role"
const PGO_ROLE_BINDING = "pgo-role-binding"
const PGO_BACKREST_ROLE = "pgo-backrest-role"
const PGO_BACKREST_ROLE_BINDING = "pgo-backrest-role-binding"

//pgo-role-binding.json
type PgoRoleBinding struct {
	TargetNamespace      string
	PgoOperatorNamespace string
}

//pgo-backrest-role.json
type PgoBackrestRole struct {
	TargetNamespace string
}

//pgo-backrest-role-binding.json
type PgoBackrestRoleBinding struct {
	TargetNamespace string
}

//pgo-role.json
type PgoRole struct {
	TargetNamespace string
}

func main() {
	apiserver.ConnectToKube()
	operatorNamespace := "pgo"
	targetNamespace := "pgouser3"
	err := apiserver.Pgo.GetConfig(apiserver.Clientset, operatorNamespace)
	if err != nil {
		log.Error("error in Pgo configuration")
		os.Exit(2)
	}

	err = CreatePGORole(apiserver.Clientset, targetNamespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}
	err = CreatePGORoleBinding(apiserver.Clientset, targetNamespace, operatorNamespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	err = CreatePGOBackrestRole(apiserver.Clientset, targetNamespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}
	err = CreatePGOBackrestRoleBinding(apiserver.Clientset, targetNamespace)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

}

func CreatePGORoleBinding(clientset *kubernetes.Clientset, targetNamespace, operatorNamespace string) error {
	//check for rolebinding existing
	_, found, _ := kubeapi.GetRoleBinding(clientset, PGO_ROLE_BINDING, targetNamespace)
	if found {
		log.Infof("rolebinding %s already exists, will not create", PGO_ROLE_BINDING)
		return nil
	}
	var buffer bytes.Buffer
	err := config.PgoRoleBindingTemplate.Execute(&buffer,
		PgoRoleBinding{
			TargetNamespace:      targetNamespace,
			PgoOperatorNamespace: operatorNamespace,
		})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(buffer.String())

	rb := v1.RoleBinding{}
	err = json.Unmarshal(buffer.Bytes(), &rb)
	if err != nil {
		log.Error("error unmarshalling " + config.PGORoleBindingPath + " json RoleBinding " + err.Error())
		return err
	}

	err = kubeapi.CreateRoleBinding(clientset, &rb, targetNamespace)
	if err != nil {
		return err
	}

	return err

}

func CreatePGOBackrestRole(clientset *kubernetes.Clientset, targetNamespace string) error {
	//check for role existing
	_, found, _ := kubeapi.GetRole(clientset, PGO_BACKREST_ROLE, targetNamespace)
	if found {
		log.Infof("role %s already exists, will not create", PGO_BACKREST_ROLE)
		return nil
	}

	var buffer bytes.Buffer
	err := config.PgoBackrestRoleTemplate.Execute(&buffer,
		PgoBackrestRole{
			TargetNamespace: targetNamespace,
		})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(buffer.String())
	r := v1.Role{}
	err = json.Unmarshal(buffer.Bytes(), &r)
	if err != nil {
		log.Error("error unmarshalling " + config.PGOBackrestRolePath + " json Role " + err.Error())
		return err
	}

	err = kubeapi.CreateRole(clientset, &r, targetNamespace)
	if err != nil {
		return err
	}

	return err
}

func CreatePGORole(clientset *kubernetes.Clientset, targetNamespace string) error {
	//check for role existing
	_, found, _ := kubeapi.GetRole(clientset, PGO_ROLE, targetNamespace)
	if found {
		log.Infof("role %s already exists, will not create", PGO_ROLE)
		return nil
	}

	var buffer bytes.Buffer
	err := config.PgoRoleTemplate.Execute(&buffer,
		PgoRole{
			TargetNamespace: targetNamespace,
		})

	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(buffer.String())
	r := v1.Role{}
	err = json.Unmarshal(buffer.Bytes(), &r)
	if err != nil {
		log.Error("error unmarshalling " + config.PGORolePath + " json Role " + err.Error())
		return err
	}

	err = kubeapi.CreateRole(clientset, &r, targetNamespace)
	if err != nil {
		return err
	}
	return err
}

func CreatePGOBackrestRoleBinding(clientset *kubernetes.Clientset, targetNamespace string) error {

	//check for rolebinding existing
	_, found, _ := kubeapi.GetRoleBinding(clientset, PGO_BACKREST_ROLE_BINDING, targetNamespace)
	if found {
		log.Infof("rolebinding %s already exists, will not create", PGO_BACKREST_ROLE_BINDING)
		return nil
	}
	var buffer bytes.Buffer
	err := config.PgoBackrestRoleBindingTemplate.Execute(&buffer,
		PgoBackrestRoleBinding{
			TargetNamespace: targetNamespace,
		})
	if err != nil {
		log.Error(err.Error() + " on " + config.PGOBackrestRoleBindingPath)
		return err
	}
	log.Info(buffer.String())

	rb := v1.RoleBinding{}
	err = json.Unmarshal(buffer.Bytes(), &rb)
	if err != nil {
		log.Error("error unmarshalling " + config.PGOBackrestRoleBindingPath + " json RoleBinding " + err.Error())
		return err
	}

	err = kubeapi.CreateRoleBinding(clientset, &rb, targetNamespace)
	if err != nil {
		return err
	}
	return err
}
