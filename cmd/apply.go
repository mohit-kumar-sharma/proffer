// Package cmd is used to created the cli utility for proffer.
/*
Copyright © 2020 mohit-kumar-sharma <flashtaken1@gmail.com>

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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"github.com/proffer/command"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	inventoryPath string

	applyLong = dedent.Dedent(`
		Apply command is used to apply the proffer configuration and distribute the cloud image
		in between multiple regions and with multiple accounts.`)

	applyExamples = dedent.Dedent(`
		$ proffer apply [flags] TEMPLATE
		$ proffer apply proffer.yml
		$ proffer apply -d proffer.yml`)

	// applyCmd represents the apply command used to apply the given configuration.
	applyCmd = &cobra.Command{
		Use:     "apply",
		Short:   "Apply proffer configuration file.",
		Long:    applyLong,
		Example: applyExamples,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("proffer config file is missing in arguments, pls pass config file to apply")
			}
			return nil
		},
		Run: applyConfig,
	}
)

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&inventoryPath, "inventory", "i", "inventory.yml", "yml file path to store proffer inventory report.")
}

// applyConfig applies the given template configuration.
func applyConfig(cmd *cobra.Command, args []string) {
	// validate template before applying
	clogger.SetPrefix("start-validation| ")
	fmt.Println()
	clogger.Info("Validating template before applying...")
	validateConfig(cmd, args)
	fmt.Println()

	clogger.SetPrefix("start-apply | ")
	clogger.Info("Applying template config...")

	if len(args) == 0 {
		log.Fatalln("Proffer Configuration file missing: Pls pass proffer config file to apply")
	}

	cfgFileAbsPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// apply template
	executeResources(cfgFileAbsPath)

	// cleanup temp files.
	if !debug {
		_ = os.Remove("output.yml")
	}
}

// executeResources applies the given resources in given configuration.
func executeResources(dsc string) {
	c, err := parseConfig(dsc)
	if err != nil {
		clogger.Fatal("Unable to parse configuration file: \n", err)
	}

	resources := command.Resources
	inventory := make([]byte, 0)
	header := []byte("---\n")
	inventory = append(inventory, header...)

	// apply resources defined in template one by one
	for _, rawResource := range c.RawResources {
		resource, ok := resources[rawResource.Type]
		if !ok {
			clogger.Fatalf("InvalidResourceType: Resource Type '%s' Not Found", rawResource.Type)
		}

		clogger.SetPrefix(rawResource.Type + " | ")
		clogger.Successf("Resource : %s  Status: Started", rawResource.Name)
		clogger.Info("")

		if err := resource.Prepare(rawResource.Config); err != nil {
			clogger.Error(err)
			clogger.Fatalf("Resource : %s  Status: Failed", rawResource.Name)
		}

		if err := resource.Run(); err != nil {
			clogger.Error(err)
			clogger.Fatalf("Resource : %s  Status: Failed", rawResource.Name)
		}

		bs, err := resource.GenerateInventory()
		if err != nil {
			clogger.Fatal(err)
		}

		inventory = append(inventory, bs...)

		clogger.Info("")
		clogger.Successf("Resource : %s  Status: Succeeded", rawResource.Name)
		fmt.Println("")
	}

	// Generate inventory report
	file, err := os.Create(inventoryPath)
	if err != nil {
		clogger.Fatal(err)
	}

	_, err = file.Write(inventory)
	if err != nil {
		clogger.Fatal(err)
	}
}
