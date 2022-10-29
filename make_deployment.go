package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// "k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
)

func main() {
	// user defined
	numReplicas := int32(10000)
	numContainers := 75
	imageName := "75tags1sig"
	deployment := createDeployment(&numReplicas, numContainers, imageName)
	bytes, err := yaml.Marshal(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func createDeployment(numReplicas *int32, numContainers int, imageName string) *appsv1.Deployment {
	containers := make([]corev1.Container, numContainers)
	for i := 1; i <= numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i),
			Image: fmt.Sprintf("%s:%v", imageName, i),
		}
	}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "perf-deployment",
			Namespace: "governed",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: numReplicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
		},
	}
}
