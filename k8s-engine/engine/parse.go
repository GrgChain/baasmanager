package engine

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"bytes"
	"log"
	"encoding/json"
)

const (
	Separator         = "---"
	ApiVersion        = "v1"
	DeploymentVersion = "apps/v1beta1"

	KindDeployment            = "Deployment"
	KindNamespace             = "Namespace"
	KindService               = "Service"
	KindPersistentVolume      = "PersistentVolume"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindStatefulSet           = "StatefulSet"
)

type ParseProject struct {
	Data [][]byte `json:"data"`
}

func (p ParseProject) Yamls2Jsons() [][]byte {
	jsons := make([][]byte, 0)
	for _, yamlBytes := range p.Data {
		yamls := bytes.Split(yamlBytes, []byte(Separator))
		for _, v := range yamls {
			if len(v) == 0 {
				continue
			}
			obj, err := yaml.ToJSON(v)
			if err != nil {
				log.Println(err.Error())
			}
			jsons = append(jsons, obj)
		}

	}
	return jsons
}

func (p ParseProject) Bytes2K8sEntities(jsonArray [][]byte) *ItemQueue {

	queue := NewItemQueue()

	for _, jsonObj := range jsonArray {
		var entity map[string]interface{}
		err := json.Unmarshal(jsonObj, &entity)
		if err != nil {
			log.Println(err.Error())
		}
		if entity["apiVersion"] == ApiVersion || entity["apiVersion"] == DeploymentVersion {
			switch entity["kind"] {
			case KindNamespace:
				namespace := &corev1.Namespace{}
				err = json.Unmarshal(jsonObj, &namespace)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(namespace)
			case KindDeployment:
				deployment := &appsv1.Deployment{}
				err = json.Unmarshal(jsonObj, &deployment)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(deployment)
			case KindService:
				service := &corev1.Service{}
				err = json.Unmarshal(jsonObj, &service)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(service)
			case KindPersistentVolume:
				pv := &corev1.PersistentVolume{}
				err = json.Unmarshal(jsonObj, &pv)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(pv)
			case KindPersistentVolumeClaim:
				pvc := &corev1.PersistentVolumeClaim{}
				err = json.Unmarshal(jsonObj, &pvc)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(pvc)
			case KindStatefulSet:
				sfs := &appsv1.StatefulSet{}
				err = json.Unmarshal(jsonObj, &sfs)
				if err != nil {
					log.Println(err.Error())
				}
				queue.Enqueue(sfs)
			}
		}
	}
	return queue
}
