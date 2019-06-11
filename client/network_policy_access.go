package client

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNetworkPolicies(clientset *kubernetes.Clientset, namespace string) (*networkingv1.NetworkPolicyList, error) {
	npi := clientset.NetworkingV1().NetworkPolicies(namespace)

	var label, field string
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	nps, err := npi.List(listOptions)
	if err != nil {
		return nil, err
	}
	return nps, nil
}
