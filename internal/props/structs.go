package props

import (
	"github.com/cantylv/thumbnail-loader/internal/entity"
	"github.com/cantylv/thumbnail-loader/services"
	"go.uber.org/zap"
)

type Load struct {
	VideoId           string
	CacheInmemoryNeed bool
	Resolutions       []int
	ServiceCluster    *services.Services
	Logger            *zap.Logger
}

func GetLoad(videoId string, cache bool, resolutions []int, cluster *services.Services, logger *zap.Logger) *Load {
	return &Load{
		VideoId:           videoId,
		CacheInmemoryNeed: cache,
		Resolutions:       resolutions,
		ServiceCluster:    cluster,
		Logger:            logger,
	}
}

type LoadDataFromServer struct {
	VideoId        string
	ServiceCluster *services.Services
	Logger         *zap.Logger
}

func GetLoadDataFromServer(videoId string, cluster *services.Services, logger *zap.Logger) *LoadDataFromServer {
	return &LoadDataFromServer{
		VideoId:        videoId,
		ServiceCluster: cluster,
		Logger:         logger,
	}
}

type SaveS3 struct {
	ImageData  map[entity.ThumbnailBody][]byte
	BucketName string
	Dir        string
	VideoId    string
	Cluster    *services.Services
	Logger     *zap.Logger
}

func GetSaveS3(imageData map[entity.ThumbnailBody][]byte, bucketName string, dir string, videoId string, cluster *services.Services, logger *zap.Logger) *SaveS3 {
	return &SaveS3{
		ImageData:  imageData,
		BucketName: bucketName,
		Dir:        dir,
		VideoId:    videoId,
		Cluster:    cluster,
		Logger:     logger,
	}
}

type DownloadS3 struct {
	ObjectName string
	Cluster    *services.Services
	Logger     *zap.Logger
}

func GetDownloadS3(objectName string, cluster *services.Services, logger *zap.Logger) *DownloadS3 {
	return &DownloadS3{
		ObjectName: objectName,
		Cluster:    cluster,
		Logger:     logger,
	}
}
