package connect

import (
	"context"
	"crypto/tls"
	"net/http"

	"connectrpc.com/connect"
	"github.com/minio/minio-go/v7"
	v1 "github.com/yoshino-s/image-downloader/gen/v1"
	"github.com/yoshino-s/image-downloader/gen/v1/v1connect"
	"github.com/yoshino-s/image-downloader/s3"
)

var _ v1connect.ImageServiceHandler = (*ImageService)(nil)

type ImageService struct {
	s3 *s3.S3
}

func NewImageService() *ImageService {
	return &ImageService{}
}

func (s *ImageService) SetS3(s3 *s3.S3) {
	s.s3 = s3
}

func (s *ImageService) Download(ctx context.Context, req *connect.Request[v1.DownloadRequest]) (*connect.Response[v1.DownloadResponse], error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: req.Msg.SkipSslVerify},
		},
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, req.Msg.Url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range req.Msg.Headers {
		r.Header.Set(k, v)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if url, err := s.s3.UploadStream(ctx, req.Msg.Key, resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: resp.Header.Get("Content-Type"),
	}); err != nil {
		return nil, err
	} else {
		return &connect.Response[v1.DownloadResponse]{
			Msg: &v1.DownloadResponse{
				Url: url.String(),
			},
		}, nil
	}

}

func (s *ImageService) Head(ctx context.Context, req *connect.Request[v1.HeadRequest]) (*connect.Response[v1.HeadResponse], error) {
	if url, err := s.s3.StatObject(ctx, req.Msg.Key, minio.StatObjectOptions{}); err != nil {
		return nil, err
	} else {
		var u string
		if url != nil {
			u = url.String()
		}

		return &connect.Response[v1.HeadResponse]{
			Msg: &v1.HeadResponse{
				Exists: url != nil,
				Url:    u,
			},
		}, nil
	}
}
