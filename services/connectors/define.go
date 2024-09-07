package connectors

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// //// CACHE //////
type ClientCache interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
	Close() error
}

var _ ClientCache = (*memcache.Client)(nil)

type EngineCache interface {
	InitClientCache(logger *zap.Logger) ClientCache
}

// //// DATABASE ///////
// ClientS3 - s3 client interface
type ClientS3 interface {
	PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	GetObject(ctx context.Context, bucketName string, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
}

var _ ClientS3 = (*minio.Client)(nil)

// EngineS3 - initialization interface for client
type EngineS3 interface {
	InitClientS3(logger *zap.Logger) ClientS3
}

type ClientDB interface {
	Conn(ctx context.Context) (*sql.Conn, error)
	Close() error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

type EngineDB interface {
	InitClientDB(logger *zap.Logger) ClientDB
}
