/*
Copyright © 2022 sharo-jef sharo.jef@gmail.com

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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type JSON map[string]interface{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cq [file...]", // 本来は "cq <cq filter> [file...]"
	Short: "csv processor like jq",
	Long:  `csv processor like jq`,
	Args: func(cmd *cobra.Command, args []string) error {
		// cq filter が実装されていないのでコメントアウト
		// if len(args) < 1 {
		// 	return errors.New("query argument is required")
		// }
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fss, _ := cmd.Flags().GetString("field-separator")
		c, _ := cmd.Flags().GetBool("compact")
		hs, _ := cmd.Flags().GetString("header")
		var header []string
		if len(hs) > 0 {
			header = strings.Split(hs, ",")
		} else {
			header = []string{}
		}

		fs, _, _, err := strconv.UnquoteChar(fss, '"')
		if err != nil {
			fmt.Printf("Failed to read field separator: %s\n", fss)
			os.Exit(2)
		}

		if len(args) == 0 {
			stdin := os.Stdin
			defer stdin.Close()

			r := csv.NewReader(stdin)
			r.FieldsPerRecord = -1
			r.Comma = fs
			fmt.Println(csvToJson(r, c, header))
		} else {
			for _, fileName := range args {
				fp, err := os.Open(fileName)
				if err != nil {
					fmt.Printf("Failed to open: %s\n", fileName)
					os.Exit(4)
				}
				defer fp.Close()

				r := csv.NewReader(fp)
				r.FieldsPerRecord = -1
				r.Comma = fs
				fmt.Println(csvToJson(r, c, header))
			}
		}
	},
}

func csvToJson(csvReader *csv.Reader, compact bool, header []string) string {
	var err error
	results := []JSON{}

	if len(header) == 0 {
		header, err = csvReader.Read()
		if err == io.EOF {
			return ""
		}
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Failed to read row:\n\n%s\n", err)
			os.Exit(5)
		}

		jsonData := make(JSON)
		for i := range row {
			if len(header) <= i {
				continue
			}
			jsonData[header[i]] = string(row[i])
		}

		results = append(results, jsonData)
	}

	var jsonBytes []byte
	if compact {
		jsonBytes, err = json.Marshal(results)
		if err != nil {
			fmt.Printf("Failed to convert to json\n\n%s\n", err)
			os.Exit(5)
		}
	} else {
		jsonBytes, err = json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Failed to convert to json\n\n%s\n", err)
			os.Exit(5)
		}
	}

	return string(jsonBytes)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cq.yaml)") // 現状 config に保存する情報がないのでコメントアウト
	rootCmd.PersistentFlags().StringP("field-separator", "F", ",", "field separator")
	rootCmd.PersistentFlags().BoolP("compact", "c", false, "compact instead of pretty-printed output")
	rootCmd.PersistentFlags().StringP("header", "H", "", "header (comma separated string)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cq" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cq")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
