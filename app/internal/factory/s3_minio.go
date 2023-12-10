package factory

import (
	"context"
	"go-sample/config"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewS3Minio(cfg *config.Config, logger mrcore.Logger) (*mrminio.ConnAdapter, error) {
	logger.Info("Create and init S3 minio connection")

	opt := mrminio.Options{
		Host:     cfg.S3.Host,
		Port:     cfg.S3.Port,
		UseSSL:   cfg.S3.UseSSL,
		User:     cfg.S3.Username,
		Password: cfg.S3.Password,
	}

	conn := mrminio.New(cfg.S3.CreateBuckets)

	if err := conn.Connect(opt); err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func RegisterS3ImageStorage(
	cfg *config.Config,
	pool *mrstorage.FileProviderPool,
	conn *mrminio.ConnAdapter,
	logger mrcore.Logger,
) error {
	storage, err := newS3MinioFileProvider(
		conn,
		cfg.FileProviders.ImageStorage.BucketName,
		logger,
	)

	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newS3MinioFileProvider(
	conn *mrminio.ConnAdapter,
	bucketName string,
	logger mrcore.Logger,
) (*mrminio.FileProvider, error) {
	logger.Info("Create and init file provider with bucket '%s'", bucketName)

	created, err := conn.InitBucket(context.Background(), bucketName)

	if err != nil {
		return nil, err
	}

	if created {
		mrcore.LogInfo("Bucket '%s' created", bucketName)
	} else {
		mrcore.LogInfo("Bucket '%s' exists, OK", bucketName)
	}

	return mrminio.NewFileProvider(conn, bucketName), nil
}
