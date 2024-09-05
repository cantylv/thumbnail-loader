package urls

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cantylv/thumbnail-loader/microservice/loader/internal/props"
	cache "github.com/cantylv/thumbnail-loader/microservice/loader/internal/repository/urls"
	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	"github.com/cantylv/thumbnail-loader/microservice/loader/utils/functions"
	"github.com/cantylv/thumbnail-loader/services"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecase interface {
	Download(ctx context.Context, args *gen.Args) (*emptypb.Empty, error)
}

type UsecaseLayer struct {
	gen.UnsafeDownloadManagerServer
	repoCacheDb    cache.Repo
	serviceCluster *services.Services
	logger         *zap.Logger
}

func NewUsecaseLayer(repo cache.Repo, serviceCluster *services.Services, logger *zap.Logger) *UsecaseLayer {
	return &UsecaseLayer{
		repoCacheDb:    repo,
		serviceCluster: serviceCluster,
		logger:         logger,
	}
}

// video resolutions for cache
var resolutions = []int{120, 320, 480, 640, 1280}

func (r *UsecaseLayer) Download(ctx context.Context, p *gen.DownloadProps) (*emptypb.Empty, error) {
	ids, err := functions.GetVideosId(p.Arguments.Data, r.logger)
	if err != nil {
		r.logger.Info(fmt.Sprintf("%v. EXAMPLE: app --cache_inmemory=false --async=true https://www.youtube.com/watch?v=6wTWF707WWE https://www.youtube.com/watch?v=5ZkdpWNtx58", err))
		return nil, nil
	}

	// create root folder for saving files
	err = os.MkdirAll(p.Flags.UploadFolder, 0755)
	if err != nil {
		r.logger.Error(fmt.Sprintf("error while creating folder: %v", err))
		return nil, nil
	}
	if p.Flags.Async {
		r.logger.Info("asynchronous loading started")
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			go func(wgOut *sync.WaitGroup) {
				p := props.GetLoad(id, p.Flags, resolutions, r.repoCacheDb, r.serviceCluster, r.logger)
				functions.Load(ctx, p)
				wgOut.Done()
			}(&wg)
		}
		wg.Wait()
	} else {
		r.logger.Info("synchronous loading started")
		for _, id := range ids {
			p := props.GetLoad(id, p.Flags, resolutions, r.repoCacheDb, r.serviceCluster, r.logger)
			functions.Load(ctx, p)
		}
	}
	r.logger.Info("thumbnails were succesful uploaded")
	return nil, nil
}
