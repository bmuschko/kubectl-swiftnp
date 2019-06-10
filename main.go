package main

import (
	"github.com/bmuschko/kubectl-swiftnp/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

var version = "undefined"

func main() {
	cmd.SetVersion(version)

	networkPolicyCmd := cmd.NewNetworkPolicyCommand(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := networkPolicyCmd.Execute(); err != nil {
		os.Exit(1)
	}
}