package factory

import (
	"context"
	"go-sample/config"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewS3Minio(cfg *config.Config, logger mrcore.Logger) (*mrminio.ConnAdapter , error) {
	logger.Info("Create and init S3 minio connection")

	opt := mrminio.Options{
		Host:	 cfg.S3.Host,
		Port:	 cfg.S3.Port,
		UseSSL:   cfg.S3.UseSSL,
		User:	 cfg.S3.Username,
		Password: cfg.S3.Password,
	}

	conn := mrminio.New()

	if err := conn.Connect(opt); err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

func NewS3MinioFileProvider(conn *mrminio.ConnAdapter, bucketName string, logger mrcore.Logger) (mrstorage.FileProviderAPI , error) {
	logger.Info("Init S3 minio bucket '%s' and create if not exists", bucketName)

	created, err := conn.InitBucket(context.Background(), bucketName, true)

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
