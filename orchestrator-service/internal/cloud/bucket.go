package cloud

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type cloud struct {
	session *session.Session
}

func New(s *session.Session) *cloud {
	return &cloud{
		session: s,
	}
}

func (c *cloud) Upload(img []byte, key string) (string, error) {
	up := s3manager.NewUploader(c.session)

	bucketName := os.Getenv("BUCKET_NAME")
	_, err := up.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(img),
		ContentType: aws.String("image/png"),
	})

	return "http://localhost:4566/" + bucketName + "/" + key, err
}
