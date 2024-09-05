package functions

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/entity"
	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/props"
	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	"github.com/cantylv/thumbnail-loader/services"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Load
// loads data from cache on hit, in another way load data from youtube servers
func Load(ctx context.Context, p *props.Load) {
	// it's for in-memory storages
	cacheHitSuccess := make([]int, 0, len(p.Resolutions))
	countSuccessHit := 0
	if p.Flags.NeedCache {
		// key == video_id + width
		// value == minio url
		for i, imgResolutionWidth := range p.Resolutions {
			item, err := p.ServiceCluster.InMemoryCacheClient.Get(fmt.Sprintf("%s%d", p.VideoId, imgResolutionWidth))
			if err != nil {
				if !errors.Is(err, memcache.ErrCacheMiss) {
					p.Logger.Info(fmt.Sprintf("internal error while cache hit: %v", err))
					return
				}
				continue
			}
			getFromCacheAndUpload(i, string(item.Value), p.Flags, p.ServiceCluster, p.Logger)
			cacheHitSuccess = append(cacheHitSuccess, imgResolutionWidth)
			countSuccessHit++
		}
	} else {
		for i, imgResolutionWidth := range p.Resolutions {
			value, err := p.RepoCache.Get(ctx, fmt.Sprintf("%s%d", p.VideoId, imgResolutionWidth))
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					p.Logger.Info(fmt.Sprintf("internal error while cache hit: %v", err))
					return
				}
				continue
			}
			getFromCacheAndUpload(i, value, p.Flags, p.ServiceCluster, p.Logger)
			cacheHitSuccess = append(cacheHitSuccess, imgResolutionWidth)
			countSuccessHit++
		}
	}
	// if no cache hit OR not everything is cached
	if countSuccessHit != len(p.Resolutions) {
		missingResolutionWidth := getMissingImageWidth(cacheHitSuccess, p.Resolutions)
		loadDataFromServerProps := props.GetLoadDataFromServer(p.VideoId, missingResolutionWidth, p.Flags, p.RepoCache, p.ServiceCluster, p.Logger)
		err := loadDataFromServer(loadDataFromServerProps)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("error while loading image from youtube server: %v", err))
			return
		}
	}
	p.Logger.Info(fmt.Sprintf("Video with id=%s was succesful uploaded", p.VideoId))
}

func getMissingImageWidth(cacheHitWidth []int, allResolutions []int) []int {
	successCacheHit := make(map[int]bool, len(cacheHitWidth))
	for _, width := range cacheHitWidth {
		successCacheHit[width] = true
	}
	res := make([]int, 0, len(allResolutions)-len(cacheHitWidth))
	for i := 0; i < len(allResolutions); i++ {
		if _, ok := successCacheHit[allResolutions[i]]; !ok {
			res = append(res, allResolutions[i])
		}
	}
	return res
}

func getFromCacheAndUpload(iteration int, value string, flags *gen.CmdFlags, cluster *services.Services, logger *zap.Logger) {
	imgUrlParts := strings.Split(value, "/")
	if iteration == 0 {
		loadPath := fmt.Sprintf("%s/%s", flags.UploadFolder, imgUrlParts[0])
		err := os.MkdirAll(loadPath, 0755)
		if err != nil {
			logger.Error(fmt.Sprintf("error while creating folder: %v", err))
			return
		}
	}
	getS3Props := props.GetDownloadS3(value, cluster, logger)
	imgData, err := getS3(getS3Props)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Info(fmt.Sprintf("video %s was received from cache", value))
	err = writeFileInDirectory(flags.UploadFolder, value, imgData)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

// loadDataFromServer
// receives json with snippet that consists of thumbnails of specific video
func loadDataFromServer(p *props.LoadDataFromServer) error {
	requestUri := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&fields=items(snippet(title,thumbnails))&part=snippet",
		p.VideoId, viper.GetString("upload.key"))
	httpResponse, err := http.Get(requestUri)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	var responseData entity.Response
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}

	loadFolder := fmt.Sprintf("%s/%s", p.Flags.UploadFolder, responseData.Items[0].Snippet.Title)
	err = os.MkdirAll(loadFolder, 0755)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("error while creating folder: %v", err))
	}
	missingThumbnails := getUncachedThumbnails(&responseData.Items[0].Snippet.Thumbnails, p.MissingResolutionWidth)

	// key - width, value - image byte
	imgUrlS3 := make(map[entity.ThumbnailBody][]byte, len(missingThumbnails))
	for _, descr := range missingThumbnails {
		imgData, err := uploadImageFromYoutube(descr, responseData.Items[0].Snippet.Title, p.Flags.UploadFolder)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("error while uploading image: %v", err.Error()))
		}
		imgUrlS3[descr] = imgData
	}

	saveS3Props := props.GetSaveS3(imgUrlS3, viper.GetString("minio.bucket_name"), responseData.Items[0].Snippet.Title, p.VideoId, p.RepoCache, p.Flags, p.ServiceCluster, p.Logger)
	saveS3AndCache(saveS3Props)
	return nil
}

func getUncachedThumbnails(thumnails *entity.ThumbnailType, resolutionWidth []int) []entity.ThumbnailBody {
	res := make([]entity.ThumbnailBody, 0, len(resolutionWidth))
	for _, width := range resolutionWidth {
		switch width {
		case int(thumnails.Default.Width):
			res = append(res, thumnails.Default)
		case int(thumnails.Medium.Width):
			res = append(res, thumnails.Medium)
		case int(thumnails.High.Width):
			res = append(res, thumnails.High)
		case int(thumnails.Standard.Width):
			res = append(res, thumnails.Standard)
		case int(thumnails.Maxres.Width):
			res = append(res, thumnails.Maxres)
		}
	}
	return res
}

// uploadImageFromYoutube
// loads image from youtube servers and saves it in directory and loads it in minio
func uploadImageFromYoutube(tBody entity.ThumbnailBody, loadFolder, uploadFolder string) ([]byte, error) {
	httpResponse, err := http.Get(tBody.Url)
	if err != nil {
		return nil, err
	}
	imageBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	loadPath := fmt.Sprintf("%s/%dx%d.jpg", loadFolder, tBody.Width, tBody.Height)
	err = writeFileInDirectory(uploadFolder, loadPath, imageBytes)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

// writeFileInDirectory
// creates file in loadFolder
func writeFileInDirectory(uploadFoler, loadPath string, data []byte) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", uploadFoler, loadPath))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// saveS3AndCache
// save image to the minio and cache it
func saveS3AndCache(p *props.SaveS3) error {
	for imgDescriptor, imgData := range p.ImageData {
		loadPath := fmt.Sprintf("%s/%dx%d.jpg", p.Dir, imgDescriptor.Width, imgDescriptor.Height)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		reader := bytes.NewReader(imgData)
		_, err := p.Cluster.MinioClient.PutObject(ctx, p.BucketName, loadPath, reader, int64(len(imgData)), minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s%d", p.VideoId, imgDescriptor.Width)
		item := memcache.Item{
			Key:        key,
			Value:      []byte(loadPath),
			Expiration: int32(p.Flags.CacheTimeout.Seconds),
		}

		// means inmemory cache
		if p.Flags.NeedCache {
			err = p.Cluster.InMemoryCacheClient.Set(&item)
			if err != nil {
				p.Logger.Info(fmt.Sprintf("error while setting value in cache: %v", err.Error()))
				continue
			}
		} else {
			err = p.RepoCache.Save(context.Background(), key, loadPath)
			if err != nil {
				p.Logger.Info(fmt.Sprintf("error while setting value in cache: %v", err.Error()))
				continue
			}
		}
	}
	return nil
}

func getS3(p *props.DownloadS3) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	obj, err := p.Cluster.MinioClient.GetObject(ctx, viper.GetString("minio.bucket_name"), p.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	return data, nil
}
