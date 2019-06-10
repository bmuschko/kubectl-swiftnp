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
		Short: "list network policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&list.namespace, "namespace", "n", "default", "the namespace used for querying for network policies")
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
		return nil
	}

	a.printNetworkPolicies(nps)
	return nil
}

func (a *networkPolicyListCmd) printNetworkPolicies(nps *v1.NetworkPolicyList) {
	if len(nps.Items) > 0 {
		table := uitable.New()
		table.AddRow("NAME", "INGRESS", "EGRESS", "POD SELECTOR")
		for _, np := range nps.Items {
			policyTypes := policyTypesToStruct(np.Spec.PolicyTypes)
			table.AddRow(np.Name, booleanIcon(policyTypes.ingress), booleanIcon(policyTypes.egress), podSelectorToString(np.Spec.PodSelector))
		}
		fmt.Fprintln(a.out, table)
	} else {
		fmt.Println("No resources found.")
	}
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

func booleanIcon(flag bool) string {
	if flag {
		return "✔"
	}

	return "✖"
}

type NetworkPolicyType struct {
	ingress bool
	egress bool
}