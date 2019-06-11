package client

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func GetPodsByLabelSelector(clientset *kubernetes.Clientset, namespace string, podSelector metav1.LabelSelector) (*corev1.PodList, error) {
	podi := clientset.CoreV1().Pods(namespace)
	var field string
	listOptions := metav1.ListOptions{
		LabelSelector: podSelectorToString(podSelector),
		FieldSelector: field,
	}
	pods, err := podi.List(listOptions)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func podSelectorToString(ls metav1.LabelSelector) string {
	var labels []string
	for key, val := range ls.MatchLabels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, val))
	}
	return strings.Join(labels, ", ")
}
