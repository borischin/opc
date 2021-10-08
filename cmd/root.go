/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"encoding/json"

	"io/ioutil"
	"sort"

	"github.com/open-policy-agent/opa/rego"
)

const REGO_DATA_ROOT = "data"

const ARG_MODULE = "module"
const ARG_QUERY_PACKAGE = "query-package"
const ARG_INPUT = "input"
const ARG_FORMAT = "format"

const (
	FORMAT_JSON     = "json"
	FORMAT_ENV_FILE = "env-file"
)

var OutputFormats = []string{
	FORMAT_JSON,
	FORMAT_ENV_FILE,
}

var queryPackage string
var moduleFile string
var inputData []string
var outputFormat string

func GetRootCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "opc",
		Short: "Open Policy Configuration",
		Long:  `This is a tool to manage your configurations via all kinds of policies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := readModule(moduleFile)
			if err != nil {
				return err
			}

			ctx := context.Background()

			// create query
			qry, err := rego.New(
				rego.Query(fmt.Sprintf("%s.%s", REGO_DATA_ROOT, queryPackage)),
				rego.Module(moduleFile, string(content)),
			).PrepareForEval(ctx)

			if err != nil {
				return err
			}

			// eval input
			rs, err := qry.Eval(ctx, rego.EvalInput(inputData))

			if err != nil {
				return err
			}

			// output map
			m := rs[0].Expressions[0].Value.(map[string]interface{})

			switch outputFormat {
			case FORMAT_ENV_FILE:
				return printEnvFileFormat(cmd, m)
			default: // json
				return printJsonFormat(cmd, m)
			}
		},
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	command.Flags().StringVarP(&moduleFile, ARG_MODULE, "m", "", "Rego file path. (required)")
	command.Flags().StringVarP(&queryPackage, ARG_QUERY_PACKAGE, "q", "main", "Query package name in rego file.")
	command.Flags().StringSliceVarP(&inputData, ARG_INPUT, "i", []string{}, "Usage: -i key1=value1 -i key2=value2 ...")
	command.Flags().StringVarP(&outputFormat, ARG_FORMAT, "f", FORMAT_JSON, fmt.Sprintf("output format(%s)", strings.Join(OutputFormats, ", ")))

	command.MarkFlagRequired(ARG_MODULE)

	return command
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(GetRootCmd().Execute())
}

func processInput(inputs []string) map[string]interface{} {
	output := map[string]interface{}{}
	for _, v := range inputs {
		parts := strings.Split(v, "=")
		output[parts[0]] = parts[1]
	}
	return output
}

func readModule(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func printEnvFileFormat(cmd *cobra.Command, input map[string]interface{}) error {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, k := range keys {
		cmd.Println(fmt.Sprintf("%s=%v", k, input[k]))
	}

	return nil
}

func printJsonFormat(cmd *cobra.Command, input map[string]interface{}) error {
	output, err := json.Marshal(input)
	if err != nil {
		return err
	}
	cmd.Println(string(output))

	return nil
}
