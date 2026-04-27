package fsx

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"
)

type AWSFS struct {
	client *s3.Client
	bucket string
}

func NewAWSFS(bucket string) (*AWSFS, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return &AWSFS{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

func (s *AWSFS) Open(name string) (fs.File, error) {
	key := strings.TrimPrefix(path.Clean(name), "/")

	result, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = result.Body.Close() }()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var size int64
	if result.ContentLength != nil {
		size = *result.ContentLength
	} else {
		size = int64(len(data))
	}

	var modTime time.Time
	if result.LastModified != nil {
		modTime = *result.LastModified
	}

	return &s3File{
		Reader:  bytes.NewReader(data),
		name:    path.Base(name),
		size:    size,
		modTime: modTime,
	}, nil
}

type s3File struct {
	*bytes.Reader
	name    string
	size    int64
	modTime time.Time
}

func (f *s3File) Close() error {
	return nil
}
func (f *s3File) Readdir(_ int) ([]os.FileInfo, error) { return nil, os.ErrInvalid }
func (f *s3File) Stat() (os.FileInfo, error)           { return f, nil }

func (f *s3File) Name() string       { return f.name }
func (f *s3File) Size() int64        { return f.size }
func (f *s3File) Mode() os.FileMode  { return 0444 }
func (f *s3File) ModTime() time.Time { return f.modTime }
func (f *s3File) IsDir() bool        { return false }
func (f *s3File) Sys() any           { return nil }
