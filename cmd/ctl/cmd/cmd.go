/*
Copyright 2020 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"

	"github.com/jetstack/cert-manager/cmd/ctl/pkg/approve"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/check"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/completion"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/convert"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/create"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/deny"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/experimental"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/inspect"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/renew"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/status"
	"github.com/jetstack/cert-manager/cmd/ctl/pkg/version"
)

func NewCertManagerCtlCommand(ctx context.Context, in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "cert-manager",
		Short: "cert-manager CLI tool to manage and configure cert-manager resources",
		Long: `
kubectl cert-manager is a CLI tool manage and configure cert-manager resources for Kubernetes`,
	}
	cmds.SetUsageTemplate(usageTemplate)

	cmds.Flags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	fakefs := flag.NewFlagSet("fake", flag.ExitOnError)
	klog.InitFlags(fakefs)
	if err := fakefs.Parse([]string{"-logtostderr=false"}); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}
	cmds.AddCommand(version.NewCmdVersion(ctx, ioStreams))
	cmds.AddCommand(convert.NewCmdConvert(ctx, ioStreams))
	cmds.AddCommand(create.NewCmdCreate(ctx, ioStreams))
	cmds.AddCommand(renew.NewCmdRenew(ctx, ioStreams))
	cmds.AddCommand(status.NewCmdStatus(ctx, ioStreams))
	cmds.AddCommand(inspect.NewCmdInspect(ctx, ioStreams))
	cmds.AddCommand(approve.NewCmdApprove(ctx, ioStreams))
	cmds.AddCommand(deny.NewCmdDeny(ctx, ioStreams))
	cmds.AddCommand(check.NewCmdCheck(ctx, ioStreams))
	cmds.AddCommand(completion.NewCmdCompletion(ctx, ioStreams))

	// Experimental features
	cmds.AddCommand(experimental.NewCmdExperimental(ctx, ioStreams))

	return cmds
}

const usageTemplate = `Usage:{{if .Runnable}}
  kubectl {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  kubectl {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "kubectl {{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
