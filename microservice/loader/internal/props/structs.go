package props

import (
	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/entity"
	cache "github.com/cantylv/thumbnail-loader/microservice/loader/internal/repository/urls"
	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	"github.com/cantylv/thumbnail-loader/services"
	"go.uber.org/zap"
)

type Load struct {
	VideoId        string
	Flags          *gen.CmdFlags
	Resolutions    []int
	RepoCache      cache.Repo
	ServiceCluster *services.Services
	Logger         *zap.Logger
}

func GetLoad(videoId string, flags *gen.CmdFlags, resolutions []int, repo cache.Repo, cluster *services.Services, logger *zap.Logger) *Load {
	return &Load{
		VideoId:        videoId,
		Flags:          flags,
		RepoCache:      repo,
		Resolutions:    resolutions,
		ServiceCluster: cluster,
		Logger:         logger,
	}
}

type LoadDataFromServer struct {
	VideoId                string
	MissingResolutionWidth []int
	Flags                  *gen.CmdFlags
	RepoCache              cache.Repo
	ServiceCluster         *services.Services
	Logger                 *zap.Logger
}

func GetLoadDataFromServer(videoId string, missingResolutionWidth []int, flags *gen.CmdFlags, repoCache cache.Repo, cluster *services.Services, logger *zap.Logger) *LoadDataFromServer {
	return &LoadDataFromServer{
		VideoId:                videoId,
		MissingResolutionWidth: missingResolutionWidth,
		Flags:                  flags,
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
	Flags      *gen.CmdFlags
	Cluster    *services.Services
	RepoCache  cache.Repo
	Logger     *zap.Logger
}

func GetSaveS3(imageData map[entity.ThumbnailBody][]byte, bucketName string, dir string, videoId string, repoCache cache.Repo, flags *gen.CmdFlags, cluster *services.Services, logger *zap.Logger) *SaveS3 {
	return &SaveS3{
		ImageData:  imageData,
		BucketName: bucketName,
		Dir:        dir,
		VideoId:    videoId,
		Flags:      flags,
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
