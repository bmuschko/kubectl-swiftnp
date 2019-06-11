package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodsByLabelSelector(clientset *kubernetes.Clientset, namespace string, podSelector *metav1.LabelSelector) (*corev1.PodList, error) {
	podi := clientset.CoreV1().Pods(namespace)
	var field string
	listOptions := metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(podSelector),
		FieldSelector: field,
	}
	pods, err := podi.List(listOptions)
	if err != nil {
		return nil, err
	}
	return pods, nil
}
