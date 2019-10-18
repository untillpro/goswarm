/*
 * Copyright (c) 2018-present unTill Pro, Ltd. and Contributors
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package swarmgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	gc "github.com/untillpro/gochips"
)

// labelsCmd represents the labels command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage swarmgo node labels",
	Long:  `Allows viewing, adding or deleting swarmgo node labels`,
	Run: loggedCmd(func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	}),
}

var labelsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List labels of the swarmgo nodes",
	Long:  `Allows viewing list of swarmgo nodes and its labels`,
	Run: loggedCmd(func(cmd *cobra.Command, args []string) {

		checkSSHAgent()
		clusterFile := unmarshalClusterYml()
		nodesList := getNodesFromYml(getWorkingDir())
		gc.ExitIfFalse(len(nodesList) > 0, "No nodes found in nodes.yml")

		var cmdline bytes.Buffer
		cmdline.WriteString("sudo docker node inspect")

		for _, node := range nodesList {
			cmdline.WriteString(" " + node.Alias)
		}

		client := getSSHClient(clusterFile)
		client.HideStdout = true // TODO: use hidden mode by default

		jsonstr := client.ExecOrExit(nodesList[0].Host, cmdline.String())
		var result []map[string]interface{}
		json.Unmarshal([]byte(jsonstr), &result)
		gc.ExitIfFalse(len(result) == len(nodesList), "Unexpected number of returned nodes")

		gc.Info(fmt.Sprintf("%-30s%-50s", "NODE", "LABELS"))
		for _, n := range result {
			spec := n["Spec"].(map[string]interface{})
			description := n["Description"].(map[string]interface{})
			labels := spec["Labels"].(map[string]interface{})
			hostName := description["Hostname"].(string)
			var labelsStr bytes.Buffer
			for k, v := range labels {
				if labelsStr.Len() > 0 {
					labelsStr.WriteString(", ")
				}
				labelsStr.WriteString(k)
				value := v.(string)
				if len(value) > 0 {
					labelsStr.WriteString("=")
					labelsStr.WriteString(value)
				}
			}
			gc.Info(fmt.Sprintf("%-30s%-50s", hostName, labelsStr.String()))
		}

	}),
}
