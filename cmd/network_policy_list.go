package cmd

import (
	"fmt"
	"github.com/bmuschko/kubectl-swiftnp/client"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type networkPolicyListCmd struct {
	out       io.Writer
	clientset *kubernetes.Clientset
	namespace string
}

func newNetworkPolicyListCommand(streams genericclioptions.IOStreams) *cobra.Command {
	list := &networkPolicyListCmd{out: streams.Out}

	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list Network Policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientset, err := client.CreateClientset()
			if err != nil {
				return err
			}
			list.clientset = clientset
			return list.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&list.namespace, "namespace", "n", "default", "the namespace used for querying for Network Policies")
	return cmd
}

func (a *networkPolicyListCmd) run() error {
	nps, err := client.CollectNetworkPolicies(a.clientset, a.namespace)
	if err != nil {
		return err
	}

	err = a.printNetworkPolicies(nps)
	if err != nil {
		return err
	}

	return nil
}

func (a *networkPolicyListCmd) printNetworkPolicies(nps *networkingv1.NetworkPolicyList) error {
	if len(nps.Items) > 0 {
		table := uitable.New()
		table.AddRow("NAME", "INGRESS", "EGRESS", "SELECTED-PODS")
		for _, np := range nps.Items {
			selectedPods, err := client.CollectPodsByLabelSelector(a.clientset, a.namespace, np.Spec.PodSelector)
			if err != nil {
				return err
			}
			podNames := podListToStringArray(selectedPods)
			policyTypes := policyTypesToStruct(np.Spec.PolicyTypes)
			table.AddRow(np.Name, booleanIcon(policyTypes.ingress), booleanIcon(policyTypes.egress), joinPodNames(podNames))
		}
		fmt.Fprintln(a.out, table)
	} else {
		fmt.Println("No resources found.")
	}

	return nil
}

func policyTypesToStruct(pts []networkingv1.PolicyType) NetworkPolicyType {
	allTypes := NetworkPolicyType{}
	for _, pt := range pts {
		if pt == networkingv1.PolicyTypeIngress {
			allTypes.ingress = true
		}
		if pt == networkingv1.PolicyTypeEgress {
			allTypes.egress = true
		}
	}
	return allTypes
}

func podListToStringArray(pods *corev1.PodList) []string {
	var selectedPods []string
	for _, p := range pods.Items {
		selectedPods = append(selectedPods, p.Name)
	}
	return selectedPods
}

func joinPodNames(podNames []string) string {
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
