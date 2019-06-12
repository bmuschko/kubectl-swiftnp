package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodsByLabelSelector(clientset *kubernetes.Clientset, namespace string, podSelector *metav1.LabelSelector) (*corev1.PodList, error) {
	podi := clientset.CoreV1().Pods(namespace)
	sel, err := massagePodSelector(podSelector)
	if err != nil {
		return nil, err
	}
	var field string
	listOptions := metav1.ListOptions{
		LabelSelector: sel,
		FieldSelector: field,
	}
	pods, err := podi.List(listOptions)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func massagePodSelector(podSelector *metav1.LabelSelector) (string, error) {
	selector, err := metav1.LabelSelectorAsSelector(podSelector)
	if err != nil {
		return "", err
	}

	l := selector.String()
	if len(l) == 0 {
		return "", nil
	}

	return metav1.FormatLabelSelector(podSelector), nil
}
