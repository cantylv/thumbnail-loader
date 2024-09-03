package functions

import (
	"fmt"
	"os"
	"sync"

	"github.com/cantylv/thumbnail-loader/internal/props"
	"github.com/cantylv/thumbnail-loader/services"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// StartEngine
// perfoms app logic
func StartEngine(cluster *services.Services, logger *zap.Logger) {
	ids, err := getVideosId(logger)
	if err != nil {
		logger.Info(fmt.Sprintf("%v. EXAMPLE: app --cache_inmemory=false --async=true https://www.youtube.com/watch?v=6wTWF707WWE https://www.youtube.com/watch?v=5ZkdpWNtx58", err))
		return
	}

	// create root folder for saving files
	err = os.MkdirAll(viper.GetString("upload_folder"), 0755)
	if err != nil {
		logger.Error(fmt.Sprintf("error while creating folder: %v", err))
		return
	}

	// video resolutions for cache
	resolutions := []int{120, 240, 360, 540, 720}

	cacheInmemoryNeed := viper.GetBool("cache_inmemory")
	asyncNeed := viper.GetBool("async")
	if asyncNeed {
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			go func(wgOut *sync.WaitGroup) {
				p := props.GetLoad(id, cacheInmemoryNeed, resolutions, cluster, logger)
				load(p)
				wgOut.Done()
			}(&wg)
		}
		wg.Wait()
	} else {
		for _, id := range ids {
			p := props.GetLoad(id, cacheInmemoryNeed, resolutions, cluster, logger)
			load(p)
		}
	}
	logger.Info("thumbnails were succesful uploaded")
}
