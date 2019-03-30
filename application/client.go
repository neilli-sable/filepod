package application

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// S3Client ...
type S3Client interface {
	PutBucketPolicy(bucketName string) error
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	DeleteAllBucketObject(bucketName string) error
	AddFile(bucketName, path string, file *os.File) error
}

type s3Client struct {
	s3         s3iface.S3API
	s3Uploader s3manageriface.UploaderAPI
}

// NewS3Client コンストラクタ
func NewS3Client(sess *session.Session) S3Client {
	s3 := s3.New(sess)
	s3Uploader := s3manager.NewUploader(sess)

	return &s3Client{
		s3:         s3,
		s3Uploader: s3Uploader,
	}
}

func (client *s3Client) CreateBucket(bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := client.s3.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Code(), aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func (client *s3Client) DeleteBucket(bucketName string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := client.s3.DeleteBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Code(), aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (client *s3Client) AddFile(bucketName, path string, file *os.File) error {
	_, err := client.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
		Body:   file,
	})
	return err
}

func (client *s3Client) DeleteAllBucketObject(bucketName string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}

	for {
		result, err := client.s3.ListObjectsV2(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println(aerr.Code(), aerr.Error())
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return err
		}

		if *result.KeyCount == 0 {
			break
		}

		deleteObjects := []*s3.ObjectIdentifier{}
		for _, con := range result.Contents {
			deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{
				Key: con.Key,
			})
		}

		input2 := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3.Delete{
				Objects: deleteObjects,
				Quiet:   aws.Bool(true),
			},
		}

		_, err = client.s3.DeleteObjects(input2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *s3Client) PutBucketPolicy(bucketName string) error {
	_, err := client.s3.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AddPerm",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::` + bucketName + `/*"
        }
    ]
}`),
	})
	return err
}
