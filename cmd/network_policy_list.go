package cmd

import (
	"fmt"
	"github.com/bmuschko/kubectl-swiftnp/collector"
	"github.com/bmuschko/kubectl-swiftnp/renderer"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
		table.AddRow("NAME", "SELECTED-PODS", "INGRESS-POLICY", "EGRESS-POLICY", "INGRESS-RULE", "EGRESS-RULE", "FROM-COUNT", "TO-COUNT")
		for _, np := range nps {
			table.AddRow(np.Name, renderPodNames(np.SelectedPodNames), renderBoolean(np.PolicyType.Ingress), renderBoolean(np.PolicyType.Egress), renderBoolean(np.IngressRule), renderBoolean(np.EgressRule), np.FromCount, np.ToCount)
		}
		fmt.Fprintln(a.out, table)
	} else {
		fmt.Println("No resources found.")
	}
}

func renderPodNames(value []string) string {
	return renderLimitedString(renderer.JoinStrings(value))
}

func renderLimitedString(value string) string {
	return renderer.LimitString(value, 30)
}

func renderBoolean(value bool) string {
	return renderer.BooleanIcon(value)
}
