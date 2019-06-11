package cmd

import (
	"fmt"
	"github.com/bmuschko/kubectl-swiftnp/collector"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
	networkPolicies, err := collector.CollectNetworkPolicies(a.namespace)
	if err != nil {
		return err
	}

	a.printNetworkPolicies(networkPolicies)
	return nil
}

func (a *networkPolicyListCmd) printNetworkPolicies(nps []collector.NetworkPolicy) {
	if len(nps) > 0 {
		table := uitable.New()
		table.AddRow("NAME", "INGRESS", "EGRESS", "SELECTED-PODS")
		for _, np := range nps {
			table.AddRow(np.Name, booleanIcon(np.PolicyType.Ingress), booleanIcon(np.PolicyType.Egress), joinPodNames(np.SelectedPodNames))
		}
		fmt.Fprintln(a.out, table)
	} else {
		fmt.Println("No resources found.")
	}
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
