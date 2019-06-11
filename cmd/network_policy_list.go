package cmd

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

type networkPolicyListCmd struct {
	out       io.Writer
	namespace string
}

func newNetworkPolicyListCommand(streams genericclioptions.IOStreams) *cobra.Command {
	list := &networkPolicyListCmd{out: streams.Out}

	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list Network Policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&list.namespace, "namespace", "n", "default", "the namespace used for querying for Network Policies")
	return cmd
}

func (a *networkPolicyListCmd) run() error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	npi := clientset.NetworkingV1().NetworkPolicies(a.namespace)

	var label, field string
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	nps, err := npi.List(listOptions)
	if err != nil {
		return err
	}

	err = a.printNetworkPolicies(clientset, nps)
	if err != nil {
		return err
	}

	return nil
}

func (a *networkPolicyListCmd) printNetworkPolicies(clientset *kubernetes.Clientset, nps *v1.NetworkPolicyList) error {
	if len(nps.Items) > 0 {
		table := uitable.New()
		table.AddRow("NAME", "INGRESS", "EGRESS", "SELECTED-PODS")
		for _, np := range nps.Items {
			selectedPods, err := collectSelectedPods(clientset, a.namespace, np.Spec.PodSelector)
			if err != nil {
				return err
			}
			policyTypes := policyTypesToStruct(np.Spec.PolicyTypes)
			table.AddRow(np.Name, booleanIcon(policyTypes.ingress), booleanIcon(policyTypes.egress), podsToString(selectedPods))
		}
		fmt.Fprintln(a.out, table)
	} else {
		fmt.Println("No resources found.")
	}

	return nil
}

func collectSelectedPods(clientset *kubernetes.Clientset, namespace string, podSelector metav1.LabelSelector) ([]string, error) {
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

	var selectedPods []string
	for _, p := range pods.Items {
		selectedPods = append(selectedPods, p.Name)
	}
	return selectedPods, nil
}

func policyTypesToStruct(pts []v1.PolicyType) NetworkPolicyType {
	allTypes := NetworkPolicyType{}
	for _, pt := range pts {
		if pt == v1.PolicyTypeIngress {
			allTypes.ingress = true
		}
		if pt == v1.PolicyTypeEgress {
			allTypes.egress = true
		}
	}
	return allTypes
}

func podSelectorToString(ls metav1.LabelSelector) string {
	var labels []string
	for key, val := range ls.MatchLabels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, val))
	}
	return strings.Join(labels, ", ")
}

func podsToString(podNames []string) string {
	return strings.Join(podNames, ", ")
}

func booleanIcon(flag bool) string {
	if flag {
		return "✔"
	}

	return "✖"
}

type NetworkPolicyType struct {
	ingress bool
	egress  bool
}
