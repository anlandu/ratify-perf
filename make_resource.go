package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"
)

const imagePrefix string = "resource"

func main() {
	// user defined
	resourceType := "job"
	namespace := "demo"
	numPods := 40
	numJobs := 40
	numReplicas := int32(1)
	numContainers := 5
	imageName := "1sigs5subjects"
	addImagePrefix := true
	registryName := "artifactstest.azurecr.io"
	switch resourceType {
	case "deployment":
		templateString, err := createDeployment(&numReplicas, numContainers, imageName, registryName)
		if err != nil {
			panic(err)
		}
		fmt.Println(templateString)
	case "pod":
		templateString, err := createPods(namespace, numPods, numContainers, imageName, registryName, addImagePrefix)
		if err != nil {
			panic(err)
		}
		fmt.Println(templateString)
	case "job":
		templateString, err := createJobs(namespace, numJobs, numContainers, imageName, registryName, addImagePrefix)
		if err != nil {
			panic(err)
		}
		fmt.Println(templateString)
	default:
		panic(fmt.Errorf("failed to find matching resource type to generate: %s", resourceType))
	}
}

func createDeployment(numReplicas *int32, numContainers int, imageName string, registryName string) (string, error) {
	containers := make([]corev1.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i+1),
			Image: fmt.Sprintf("%s/%s:%v", registryName, imageName, i+1),
		}
	}
	template := &appsv1.Deployment{
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

	bytes, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to create deployment: %v", err)
	}
	return string(bytes), nil
}

func createPod(namespace string, podName string, numContainers int, imageName string, registryName string) *corev1.Pod {
	containers := make([]corev1.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i+1),
			Image: fmt.Sprintf("%s/%s:%v", registryName, imageName, i+1),
		}
	}
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}
}

func createPods(namespace string, numPods int, numContainers int, imageName string, registryName string, addImagePrefix bool) (string, error) {

	var podSpec *corev1.Pod
	retTemplate := ""
	scopedImageName := imageName
	for i := 1; i <= numPods; i++ {
		if addImagePrefix {
			scopedImageName = fmt.Sprintf("%s-%d-%s", imagePrefix, i, imageName)
		}
		podSpec = createPod(namespace, fmt.Sprintf("test-pod-%d", i), numContainers, scopedImageName, registryName)
		bytes, err := yaml.Marshal(podSpec)
		if err != nil {
			return "", fmt.Errorf("failed to create pod %d: %v", i, err)
		}
		retTemplate = fmt.Sprintf("%s%s\n---\n", retTemplate, string(bytes))
	}
	return retTemplate, nil
}

func createJob(namespace string, jobName string, numContainers int, imageName string, registryName string) *batchv1.Job {
	containers := make([]corev1.Container, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i] = corev1.Container{
			Name:  fmt.Sprintf("%s%v", imageName, i+1),
			Image: fmt.Sprintf("%s/%s:%v", registryName, imageName, i+1),
		}
	}
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", jobName, "pod"),
					Namespace: namespace,
				},
				Spec: corev1.PodSpec{
					Containers:    containers,
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}
}

func createJobs(namespace string, numJobs int, numContainers int, imageName string, registryName string, addImagePrefix bool) (string, error) {

	var jobSpec *batchv1.Job
	retTemplate := ""
	scopedImageName := imageName
	for i := 1; i <= numJobs; i++ {
		if addImagePrefix {
			scopedImageName = fmt.Sprintf("%s-%d-%s", imagePrefix, i, imageName)
		}
		jobSpec = createJob(namespace, fmt.Sprintf("test-job-%d", i), numContainers, scopedImageName, registryName)
		bytes, err := yaml.Marshal(jobSpec)
		if err != nil {
			return "", fmt.Errorf("failed to create job %d: %v", i, err)
		}
		retTemplate = fmt.Sprintf("%s%s\n---\n", retTemplate, string(bytes))
	}
	return retTemplate, nil
}
