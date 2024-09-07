package functions

import (
	"errors"
	"testing"

	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/props"
	"github.com/cantylv/thumbnail-loader/services"
	"github.com/cantylv/thumbnail-loader/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetS3InvalidBucketName(t *testing.T) {
	ctrl := gomock.NewController(t)
	objectName := "120x120.jpg"
	errBucketName := errors.New("invalid bucket name")

	dbCacheClientMock := mocks.NewMockClientDB(ctrl)
	inMemoryCacheClientMock := mocks.NewMockClientCache(ctrl)
	minioClientMock := mocks.NewMockClientS3(ctrl)
	minioClientMock.
		EXPECT().
		GetObject(gomock.Any(), "", objectName, minio.GetObjectOptions{}).
		Return(nil, errBucketName)
	p := props.GetDownloadS3(objectName, &services.Services{
		DBCacheClient:       dbCacheClientMock,
		InMemoryCacheClient: inMemoryCacheClientMock,
		MinioClient:         minioClientMock,
	}, zap.Must(zap.NewProduction()))
	imgData, err := getS3(p)
	assert.Equal(t, errBucketName, err)
	assert.Nil(t, imgData)
}

func TestGetS3(t *testing.T) {
	ctrl := gomock.NewController(t)
	objectName := "120x120.jpg"

	dbCacheClientMock := mocks.NewMockClientDB(ctrl)
	inMemoryCacheClientMock := mocks.NewMockClientCache(ctrl)
	minioClientMock := mocks.NewMockClientS3(ctrl)
	// expected := []byte{123, 23, 23, 231}
	minioClientMock.
		EXPECT().
		GetObject(gomock.Any(), "images", objectName, minio.GetObjectOptions{}).
		Return(&minio.Object{}, nil)
	p := props.GetDownloadS3(objectName, &services.Services{
		DBCacheClient:       dbCacheClientMock,
		InMemoryCacheClient: inMemoryCacheClientMock,
		MinioClient:         minioClientMock,
	}, zap.Must(zap.NewProduction()))
	imgData, err := getS3(p)
	assert.NotNil(t, err)
	assert.Equal(t, nil, imgData)
}
