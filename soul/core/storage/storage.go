package storage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	storageEnabled bool
	client         *minio.Client
	bucketName     string
}

func NewS3(storageEnabled bool, host string, accessKey string, secretKey string, bucketName string) (*Storage, error) {
	if !storageEnabled {
		return &Storage{storageEnabled: false}, nil
	}
	minioClient, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return &Storage{}, errors.Wrap(err, "[s3] failed to connect")
	}

	log.Info().Msg("[s3] minio is set up")
	return &Storage{client: minioClient, bucketName: bucketName}, nil
}

func (s *Storage) Upload(ctx context.Context, filename string, file []byte) error {
	if !s.storageEnabled {
		log.Info().Msgf("[s3] not upload file %s. storage disabled", filename)
		return nil
	}
	r := bytes.NewReader(file)
	info, err := s.client.PutObject(ctx, s.bucketName, filename, r, int64(len(file)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return errors.Wrapf(err, "[s3] failed to put object %s", filename)
	}
	log.Info().Msgf("[s3] file put success: %+v", info)
	return nil
}
