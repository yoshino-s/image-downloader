package s3

import (
	"context"
	"io"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
)

type S3 struct {
	*application.EmptyApplication
	config config
	client *minio.Client
}

func New() *S3 {
	return &S3{
		EmptyApplication: application.NewEmptyApplication(),
		config:           config{},
	}
}

func (s *S3) Configuration() configuration.Configuration {
	return &s.config
}

func (s *S3) Setup(context.Context) {
	client, err := minio.New(s.config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.config.AccessKeyID, s.config.SecretAccessKey, ""),
		Secure: !s.config.Insecure,
		Region: s.config.Region,
	})
	if err != nil {
		panic(err)
	}
	s.client = client
}

func (s *S3) Client() *minio.Client {
	return s.client
}

func (s *S3) UploadStream(ctx context.Context, key string, reader io.Reader, size int64, options minio.PutObjectOptions) (*url.URL, error) {
	info, err := s.client.PutObject(ctx, s.config.Bucket, key, reader, size, options)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("uploaded", zap.Any("info", info))

	return s.GetUrl(ctx, s.config.Bucket, key)
}

func (s *S3) StatObject(ctx context.Context, key string, opts minio.StatObjectOptions) (*url.URL, error) {
	if _, err := s.client.StatObject(ctx, s.config.Bucket, key, opts); err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, err
	}
	return s.GetUrl(ctx, s.config.Bucket, key)
}

func (s *S3) GetUrl(ctx context.Context, bucket, key string) (*url.URL, error) {
	u, err := s.client.PresignedGetObject(ctx, bucket, key, s.config.PresignedGetObjectExpires, url.Values{})
	if s.config.Public {
		u.RawQuery = ""
	}
	return u, err
}
