/*

Copyright 2017 Google Inc.

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
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	pb "github.com/google/vmregistry/api"
)

var (
	outputJSON bool
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initCredStoreSession()

		ctx, err := vmregistryContext(context.Background())
		if err != nil {
			glog.Fatalf("failed to acquire a client vmregistry context: %v", err)
		}

		client, err := newClient()
		if err != nil {
			glog.Fatalf("failed to create a client: %v", err)
		}

		repl, err := client.List(ctx, &pb.ListVMRequest{})
		if err != nil {
			glog.Fatalf("failed to get list of VMs: %v", err)
		}

		if outputJSON {
			b, _ := json.Marshal(repl)
			fmt.Println(string(b))
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "MAC", "IP"})

		for _, vm := range repl.Vms {
			table.Append([]string{vm.Name, vm.Mac, vm.Ip})
		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	lsCmd.Flags().BoolVar(&outputJSON, "json", false, "Output in JSON")
}
