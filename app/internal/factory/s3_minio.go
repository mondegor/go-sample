package factory

import (
	"context"

	"go-sample/config"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewS3Minio - comment func.
func NewS3Minio(ctx context.Context, cfg config.Config) (*mrminio.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init S3 minio connection")

	opts := mrminio.Options{
		Host:     cfg.S3.Host,
		Port:     cfg.S3.Port,
		UseSSL:   cfg.S3.UseSSL,
		User:     cfg.S3.Username,
		Password: cfg.S3.Password,
	}

	conn := mrminio.New(
		cfg.S3.CreateBuckets,
		mrlib.NewMimeTypeList(cfg.MimeTypes), // TODO: можно вынести в общую переменную
	)

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

// RegisterS3ImageStorage - comment func.
func RegisterS3ImageStorage(
	ctx context.Context,
	cfg config.Config,
	pool *mrstorage.FileProviderPool,
	conn *mrminio.ConnAdapter,
) error {
	storage, err := newS3MinioFileProvider(
		ctx,
		conn,
		cfg.FileProviders.ImageStorage.BucketName,
	)
	if err != nil {
		return err
	}

	return pool.Register(cfg.FileProviders.ImageStorage.Name, storage)
}

func newS3MinioFileProvider(
	ctx context.Context,
	conn *mrminio.ConnAdapter,
	bucketName string,
) (*mrminio.FileProvider, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msgf("Create and init file provider with bucket '%s'", bucketName)

	created, err := conn.InitBucket(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if created {
		logger.Debug().Msgf("Bucket '%s' created", bucketName)
	} else {
		logger.Debug().Msgf("Bucket '%s' exists, OK", bucketName)
	}

	return mrminio.NewFileProvider(conn, bucketName), nil
}
