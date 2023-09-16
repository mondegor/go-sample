package factory

import (
    "context"
    "go-sample/config"

    "github.com/minio/minio-go/v7"
    "github.com/mondegor/go-storage/mrminio"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewS3Minio(cfg *config.Config, logger mrcore.Logger) (mrstorage.FileProvider, error) {
    logger.Info("Create and init S3 minio connection")

    opt := mrminio.Options{
        Host: cfg.S3.Host,
        Port: cfg.S3.Port,
        UseSSL: cfg.S3.UseSSL,
        User: cfg.S3.Username,
        Password: cfg.S3.Password,
    }

    conn := mrminio.New(cfg.S3.BacketName)
    err := conn.Connect(opt)

    if err != nil {
        return nil, err
    }

    err = conn.Ping(context.Background())

    if err != nil {
        return nil, err
    }

    exists, err := conn.Cli().BucketExists(context.Background(), cfg.S3.BacketName)

    if err != nil {
        return nil, err
    }

    if exists {
        logger.Info("Backet '%s' exists, OK", cfg.S3.BacketName)
    } else {
        err = conn.Cli().MakeBucket(
            context.Background(),
            cfg.S3.BacketName,
            minio.MakeBucketOptions{}, // "ru-central1"
        )

        if err != nil {
            return nil, err
        }

        logger.Info("Backet '%s' created", cfg.S3.BacketName)
    }

    return conn, nil
}
