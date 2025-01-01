package s3Service

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3ServiceImpl struct {
	S3Client *s3.Client
	Bucket   string
}

func (s *S3ServiceImpl) UploadToS3(ctx context.Context, objectKey, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error in opening file %s.Error is %v\n", fileName, err)
		return err
	}
	defer file.Close()
	_, err = s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Printf("Error in uploading file %s to S3.Error is %v\n", fileName, err)
		return err
	} else {
		err = s3.NewObjectExistsWaiter(s.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(s.Bucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
			return err
		}
	}
	return nil
}

func (s *S3ServiceImpl) DownloadFromS3(ctx context.Context, objectKey string, fileName string) error {

	result, err := s.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", objectKey, s.Bucket)
			return err
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", s.Bucket, objectKey, err)
			return err
		}
	}

	defer result.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Error in reading body of object %s from S3.Error is %v\n", objectKey, err)
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		log.Printf("Error in writing body of object %s from S3.Error is %v\n", objectKey, err)
		return err
	}
	return nil
}
