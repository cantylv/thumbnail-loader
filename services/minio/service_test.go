//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package minio

import (
	"context"
	"errors"
	"testing"

	"github.com/cantylv/thumbnail-loader/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func TestMakeBucketEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("panic was recovered in the test")
		}
	}()
	ctrl := gomock.NewController(t)
	errBucket := errors.New("bucket name is empty")
	bucketName := ""

	clientMock := mocks.NewMockClientS3(ctrl)
	clientMock.
		EXPECT().
		MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}).
		Return(errBucket)

	clientMock.
		EXPECT().
		BucketExists(context.Background(), bucketName).
		Return(false, errBucket)

	makeBucket(clientMock, bucketName, context.Background(), zap.Must(zap.NewProduction()))
}

func TestMakeBucketAlreadyExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	errBucket := errors.New("bucket name is already exist")
	bucketName := "images"

	clientMock := mocks.NewMockClientS3(ctrl)
	clientMock.
		EXPECT().
		MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}).
		Return(errBucket)

	clientMock.
		EXPECT().
		BucketExists(context.Background(), bucketName).
		Return(true, nil)

	makeBucket(clientMock, bucketName, context.Background(), zap.Must(zap.NewProduction()))
}