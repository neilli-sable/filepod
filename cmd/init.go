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
	"bufio"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/neilli-sable/filepod/application"
	"github.com/spf13/cobra"
)

var profile string
var region string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize directory",
	Long: `Initialize directory.
Create filepod.json and bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make Session
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", profile),
			Region:      aws.String(region),
		}))

		fmt.Print("BucketName? ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		bucketName := scanner.Text()

		file, err := os.Create("filepod.json")
		if err != nil {
			panic(err)
		}

		client := application.NewS3Client(sess)
		err = client.CreateBucket(bucketName)
		if err != nil {
			panic(err)
		}

		file.WriteString(`{
    "bucketName": "` + bucketName + `"
}`)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().StringVar(&profile, "profile", "default", "Use a specific profile from your credential file.")
	initCmd.PersistentFlags().StringVar(&region, "region", "", "The region to use.")
}
