package s3Service

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3ServiceImpl struct {
	S3Client *s3.Client
	Bucket   string
}

func (s *S3ServiceImpl) UploadToS3() error {
	return nil
}

func (s *S3ServiceImpl) DownloadFromS3() error {
	return nil
}
