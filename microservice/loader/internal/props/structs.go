package props

import (
	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/entity"
	cache "github.com/cantylv/thumbnail-loader/microservice/loader/internal/repository/urls"
	"github.com/cantylv/thumbnail-loader/services"
	"go.uber.org/zap"
)

type Load struct {
	VideoId           string
	CacheInmemoryNeed bool
	Resolutions       []int
	RepoCache         cache.Repo
	ServiceCluster    *services.Services
	Logger            *zap.Logger
}

func GetLoad(videoId string, cache bool, resolutions []int, repo cache.Repo, cluster *services.Services, logger *zap.Logger) *Load {
	return &Load{
		VideoId:           videoId,
		CacheInmemoryNeed: cache,
		RepoCache:         repo,
		Resolutions:       resolutions,
		ServiceCluster:    cluster,
		Logger:            logger,
	}
}

type LoadDataFromServer struct {
	VideoId                string
	MissingResolutionWidth []int
	RepoCache              cache.Repo
	ServiceCluster         *services.Services
	Logger                 *zap.Logger
}

func GetLoadDataFromServer(videoId string, missingResolutionWidth []int, repoCache cache.Repo, cluster *services.Services, logger *zap.Logger) *LoadDataFromServer {
	return &LoadDataFromServer{
		VideoId:                videoId,
		MissingResolutionWidth: missingResolutionWidth,
		ServiceCluster:         cluster,
		RepoCache:              repoCache,
		Logger:                 logger,
	}
}

type SaveS3 struct {
	ImageData  map[entity.ThumbnailBody][]byte
	BucketName string
	Dir        string
	VideoId    string
	Cluster    *services.Services
	RepoCache  cache.Repo
	Logger     *zap.Logger
}

func GetSaveS3(imageData map[entity.ThumbnailBody][]byte, bucketName string, dir string, videoId string, repoCache cache.Repo, cluster *services.Services, logger *zap.Logger) *SaveS3 {
	return &SaveS3{
		ImageData:  imageData,
		BucketName: bucketName,
		Dir:        dir,
		VideoId:    videoId,
		Cluster:    cluster,
		RepoCache:  repoCache,
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
