package functions

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/cantylv/thumbnail-loader/internal/entity"
	"github.com/cantylv/thumbnail-loader/internal/props"
	"github.com/mailru/easyjson"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// load
// loads data from cache on hit, in another way load data from youtube servers
func load(p *props.Load) {
	if p.CacheInmemoryNeed {
		// key == video_id + width
		// value == minio url
		for _, videoResolution := range p.Resolutions {
			item, err := p.ServiceCluster.InMemoryCacheClient.Get(fmt.Sprintf("%s%d", p.VideoId, videoResolution))
			if err != nil {
				if !errors.Is(err, memcache.ErrCacheMiss) {
					p.Logger.Info(fmt.Sprintf("internal error while cache hit: %v", err))
					return
				}
				continue
			}
			// loadpath == title/resolution.jpg
			loadPath := string(item.Value)
			getS3Props := props.GetDownloadS3(loadPath, p.ServiceCluster, p.Logger)
			imgData, err := getS3(getS3Props)
			if err != nil {
				p.Logger.Error(err.Error())
				return
			}
			err = writeFileInDirectory(loadPath, imgData)
			if err != nil {
				p.Logger.Error(err.Error())
				return
			}
		}
	}
	// if no cache hit
	loadDataFromServerProps := props.GetLoadDataFromServer(p.VideoId, p.ServiceCluster, p.Logger)
	err := loadDataFromServer(loadDataFromServerProps)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("error while loading image from youtube server: %v", err))
		return
	}
	p.Logger.Info(fmt.Sprintf("Video with id=%s was succesful uploaded", p.VideoId))
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
	err = easyjson.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}

	loadFolder := fmt.Sprintf("%s/%s", viper.GetString("upload_folder"), responseData.Items[0].Snippet.Title)
	err = os.MkdirAll(loadFolder, 0755)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("error while creating folder: %v", err))
	}
	snippetThumbnails := responseData.Items[0].Snippet.Thumbnails
	imagesDescriptor := []entity.ThumbnailBody{
		snippetThumbnails.Default,
		snippetThumbnails.Medium,
		snippetThumbnails.High,
		snippetThumbnails.Standard,
		snippetThumbnails.Maxres,
	}

	// key - width, value - image byte
	imgUrlS3 := make(map[entity.ThumbnailBody][]byte, len(imagesDescriptor))
	for _, descr := range imagesDescriptor {
		imgData, err := uploadImageFromYoutube(descr, responseData.Items[0].Snippet.Title)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("error while uploading image: %v", err.Error()))
		}
		imgUrlS3[descr] = imgData
	}

	saveS3Props := props.GetSaveS3(imgUrlS3, viper.GetString("minio.bucket_name"), responseData.Items[0].Snippet.Title, p.VideoId, p.ServiceCluster, p.Logger)
	saveS3(saveS3Props)
	return nil
}

// uploadImageFromYoutube
// loads image from youtube servers and saves it in directory and loads it in minio
func uploadImageFromYoutube(tBody entity.ThumbnailBody, loadFolder string) ([]byte, error) {
	httpResponse, err := http.Get(tBody.Url)
	if err != nil {
		return nil, err
	}
	imageBytes, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	loadPath := fmt.Sprintf("%s/%dx%d.jpg", loadFolder, tBody.Width, tBody.Height)
	err = writeFileInDirectory(loadPath, imageBytes)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

// writeFileInDirectory
// creates file in loadFolder
func writeFileInDirectory(loadPath string, data []byte) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", viper.GetString("upload_folder"), loadPath))
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

// saveS3
// save image to the minio and cache it
func saveS3(p *props.SaveS3) error {
	for imgDescriptor, imgData := range p.ImageData {
		loadPath := fmt.Sprintf("%s/%dx%d.jpg", p.Dir, imgDescriptor.Width, imgDescriptor.Height)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		reader := bytes.NewReader(imgData)
		_, err := p.Cluster.MinioClient.PutObject(ctx, p.BucketName, loadPath, reader, int64(len(imgData)), minio.PutObjectOptions{})
		if err != nil {
			return err
		}
		item := memcache.Item{
			Key:        fmt.Sprintf("%s%d", p.VideoId, imgDescriptor.Width),
			Value:      []byte(loadPath),
			Expiration: 600,
		}
		if viper.GetBool("cache_inmemory") {
			err = p.Cluster.InMemoryCacheClient.Set(&item)
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
