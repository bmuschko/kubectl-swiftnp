package collector

import (
	"github.com/bmuschko/kubectl-swiftnp/client"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

type NetworkPolicy struct {
	Name             string
	PolicyType       NetworkPolicyType
	SelectedPodNames []string
}

type NetworkPolicyType struct {
	Ingress bool
	Egress  bool
}

func CollectNetworkPolicies(namespace string) ([]NetworkPolicy, error) {
	var networkPolicies []NetworkPolicy

	clientset, err := client.CreateClientset()
	if err != nil {
		return nil, err
	}

	nps, err := client.GetNetworkPolicies(clientset, namespace)
	if err != nil {
		return nil, err
	}

	for _, np := range nps.Items {
		networkPolicy := NetworkPolicy{Name: np.Name, PolicyType: policyTypesToStruct(np.Spec.PolicyTypes)}
		selectedPods, err := client.GetPodsByLabelSelector(clientset, namespace, np.Spec.PodSelector)
		if err != nil {
			return nil, err
		}
		networkPolicy.SelectedPodNames = podListToString(selectedPods)
		networkPolicies = append(networkPolicies, networkPolicy)
	}

	return networkPolicies, nil
}

func policyTypesToStruct(pts []networkingv1.PolicyType) NetworkPolicyType {
	allTypes := NetworkPolicyType{}
	for _, pt := range pts {
		if pt == networkingv1.PolicyTypeIngress {
			allTypes.Ingress = true
		}
		if pt == networkingv1.PolicyTypeEgress {
			allTypes.Egress = true
		}
	}
	return allTypes
}

func podListToString(pods *corev1.PodList) []string {
	var selectedPods []string
	for _, p := range pods.Items {
		selectedPods = append(selectedPods, p.Name)
	}
	return selectedPods
}
