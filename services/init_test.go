//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package services

import (
	"testing"

	"github.com/cantylv/thumbnail-loader/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func InitServices(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("panic was intercepted: %v", r)
		}
	}()
	ctrl := gomock.NewController(t)
	logger := zap.Must(zap.NewProduction())

	inMemoryClientEngineMock := mocks.NewMockEngineCache(ctrl)
	inMemoryClientMock := mocks.NewMockClientCache(ctrl)
	inMemoryClientEngineMock.EXPECT().InitClientCache(logger).Return(inMemoryClientMock)

	dbCacheClientEngineMock := mocks.NewMockEngineDB(ctrl)
	dbCacheClientMock := mocks.NewMockClientDB(ctrl)
	dbCacheClientEngineMock.EXPECT().InitClientDB(logger).Return(dbCacheClientMock)

	minioClientEngineMock := mocks.NewMockEngineS3(ctrl)
	minioClientMock := mocks.NewMockClientDB(ctrl)
	minioClientEngineMock.EXPECT().InitClientS3(logger).Return(minioClientMock)

	cluster := Init(logger, inMemoryClientEngineMock, dbCacheClientEngineMock, minioClientEngineMock)
	assert.NotNil(t, cluster.DBCacheClient, "here must be success")
	assert.NotNil(t, cluster.InMemoryCacheClient, "here must be success")
	assert.NotNil(t, cluster.MinioClient, "here must be success")
}
