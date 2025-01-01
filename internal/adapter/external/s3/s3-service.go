package s3Service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service interface {
	UploadToS3(ctx context.Context, objectKey, fileName string) error
	DownloadFromS3(ctx context.Context, objectKey string, fileName string) error
}

func NewS3ServiceImpl() (*S3ServiceImpl, error) {
	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, "")))

	if err != nil {
		log.Println("Error in Setup for S3 ", err)
		return nil, err

	}
	client := s3.NewFromConfig(cfg)
	//timeStamp := time.Now().Format("2006-01-02 15:04:05")
	bucketName := fmt.Sprintf("upload-file-%s", AWS_REGION)

	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(AWS_REGION),
		},
	})

	s3Config := &S3ServiceImpl{
		S3Client: client,
		Bucket:   bucketName,
	}

	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists

		if errors.As(err, &owned) {
			log.Println("Bucket already created by me")
			return s3Config, nil
		} else if errors.As(err, &exists) {
			log.Println("Bucket already exists")
			return nil, exists

		}
	} else {
		err = s3.NewBucketExistsWaiter(client).Wait(
			context.Background(), &s3.HeadBucketInput{Bucket: aws.String(bucketName)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", bucketName)
			return nil, err
		}
	}
	log.Println("Bucket created successfully")

	return s3Config, nil
}
