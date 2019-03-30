// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/neilli-sable/filepod/application"
	"github.com/neilli-sable/filepod/model"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push files",
	Long:  `push files.`,
	Run: func(cmd *cobra.Command, args []string) {
		paths := dirwalk("./")

		// Make Session
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", profile),
			Region:      aws.String(region),
		}))

		var setting model.FilePodJSON
		file, err := ioutil.ReadFile("./filepod.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &setting)
		bucketName := setting.BucketName

		client := application.NewS3Client(sess)

		for _, path := range paths {
			file, err := os.Open(path)

			err = client.AddFile(bucketName, path, file)
			if err != nil {
				panic(err)
			}

		}
	},
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.PersistentFlags().StringVar(&profile, "profile", "default", "Use a specific profile from your credential file.")
	pushCmd.PersistentFlags().StringVar(&region, "region", "", "The region to use.")
}
