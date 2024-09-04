package functions

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/cantylv/thumbnail-loader/microservice/loader/proto/gen"
	e "github.com/cantylv/thumbnail-loader/microservice/loader/utils/myerrors"
	"go.uber.org/zap"
)

var (
	hasDomain = `^[a-z/:]*[www.]?youtube.com+`
)

// StartEngine
// makes grpc requests
func StartEngine(client gen.DownloadManagerClient, logger *zap.Logger) {
	// get cmd args
	cmdArgs, err := getCmdArgs(logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*len(cmdArgs))*time.Second)
	defer cancel()
	_, err = client.Download(ctx, &gen.Args{Data: cmdArgs})
	if err != nil {
		logger.Error(err.Error())
	}
}

func getCmdArgs(logger *zap.Logger) ([]string, error) {
	parseUris := make([]string, 0, len(os.Args)-1) // len(os.Args) >= 1
	for _, arg := range os.Args {
		isUrl, err := isYoutubeUrl(arg)
		if err != nil {
			logger.Warn(fmt.Sprintf("error while check programm argument: %v", err))
			continue
		}
		if isUrl {
			parseUris = append(parseUris, arg)
		}
	}
	if len(parseUris) == 0 {
		return nil, e.ErrIncorrectLinks
	}
	return parseUris, nil
}

// isYoutubeUrl
// checks that uri has youtube domain
func isYoutubeUrl(input string) (bool, error) {
	matched, err := regexp.MatchString(hasDomain, input)
	if err != nil {
		return false, nil
	}
	return matched, nil
}
