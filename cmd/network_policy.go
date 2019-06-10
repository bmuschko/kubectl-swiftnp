package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func NewNetworkPolicyCommand(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "swiftnp",
		Short:        "Retrieves information about Network Policies",
	}

	cmd.AddCommand(newNetworkPolicyListCommand(streams))
	cmd.AddCommand(newVersionCommand(streams.Out))
	return cmd
}