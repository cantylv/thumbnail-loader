package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *minio.Client {
	minioEndpoint := viper.GetString("minio.host") + ":" + viper.GetString("minio.port")
	client, err := minio.New(minioEndpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			viper.GetString("minio.access_key"),
			viper.GetString("minio.secret_access_key"),
			""),
		Secure: viper.GetBool("minio.use_ssl"),
	})
	if err != nil {
		logger.Fatal("Failed to connect to Minio", zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	makeBucket(client, viper.GetString("minio.bucket_name"), ctx, logger)
	logger.Info("Minio connected successfully")
	return client
}

func makeBucket(client *minio.Client, bucket string, ctx context.Context, logger *zap.Logger) {
	err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		isExist, err := client.BucketExists(ctx, bucket)
		if err == nil && isExist {
			logger.Info(fmt.Sprintf("A bucket with a name %s already exists", bucket))
			return
		} else {
			logger.Fatal(fmt.Sprintf("Creating a bucket with a name %s was failed", bucket), zap.String("error", err.Error()))
		}
	}
	logger.Info(fmt.Sprintf("A bucket with a name %s was created successfully", bucket))
}
