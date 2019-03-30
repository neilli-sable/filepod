package application

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/neilli-sable/filepod/mock"
)

func TestCreateBucket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s3Mock := mock.NewMockS3API(ctrl)
	s3Mock.EXPECT().CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("新しいバケット")}).Return(nil, nil)

	client := s3Client{
		s3: s3Mock,
	}

	err := client.CreateBucket("新しいバケット")
	if err != nil {
		t.Fatal(err)
	}
}
