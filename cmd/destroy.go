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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/neilli-sable/filepod/application"
	"github.com/neilli-sable/filepod/model"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Delete bucket and filepod.json",
	Long:  `Delete bucket and filepod.json`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make Session
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", profile),
			Region:      aws.String(region),
		}))

		file, err := ioutil.ReadFile("./filepod.json")
		if err != nil {
			panic(err)
		}

		var setting model.FilePodJSON

		json.Unmarshal(file, &setting)
		bucketName := setting.BucketName

		client := application.NewS3Client(sess)
		err = client.DeleteBucket(bucketName)
		if err != nil {
			panic(err)
		}

		client.DeleteAllBucketObject(bucketName)
		if err != nil {
			panic(err)
		}

		err = os.Remove("filepod.json")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.PersistentFlags().StringVar(&profile, "profile", "default", "Use a specific profile from your credential file.")
	destroyCmd.PersistentFlags().StringVar(&region, "region", "", "The region to use.")
}
