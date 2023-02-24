package cloud

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	port        = "http://localhost:4566/"
	FilePath    = "../../tmp/"
	contentType = "image/jpeg"
)

type cloud struct {
	session *session.Session
}

func New(s *session.Session) *cloud {
	return &cloud{
		session: s,
	}
}

func (c *cloud) Download(key string) (string, error) {
	up := s3manager.NewDownloader(c.session)
	file, err := ioutil.TempFile(FilePath, "img-*.jpeg")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bucketName := os.Getenv("BUCKET_NAME")

	_, err = up.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})

	return file.Name(), err
}

func (c *cloud) Upload(img []byte, key string) (string, error) {
	up := s3manager.NewUploader(c.session)

	bucketName := os.Getenv("BUCKET_NAME2")
	_, err := up.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(img),
		ContentType: aws.String(contentType),
	})

	return port + bucketName + "/" + key, err
}
