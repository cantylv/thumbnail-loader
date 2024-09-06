//go:generate go test ./... -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html
package functions

import (
	"os"
	"testing"

	"github.com/cantylv/thumbnail-loader/microservice/loader/mocks"
	e "github.com/cantylv/thumbnail-loader/microservice/loader/utils/myerrors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIsYoutubeUrlInvalidUrls(t *testing.T) {
	isUrlYoutube, err := isYoutubeUrl("https://bmstu.ru/iu5/education")
	assert.Nil(t, err, "other valid URLs should not throw an error")
	assert.Equal(t, false, isUrlYoutube, "domen of url must be [www.]youtube.com")

	isUrlYoutube, err = isYoutubeUrl("https://")
	assert.Nil(t, err, "only schema should not throw an error")
	assert.Equal(t, false, isUrlYoutube, "no youtube domain")

	isUrlYoutube, err = isYoutubeUrl("")
	assert.Nil(t, err, "empty string should not throw an error")
	assert.Equal(t, false, isUrlYoutube, "empty string is not youtube url")
}

func TestIsYoutubeUrlValidUrls(t *testing.T) {
	isUrlYoutube, err := isYoutubeUrl("https://www.youtube.com/watch?v=q-cE4ziUBMo")
	assert.Nil(t, err, "valid youtube url should not throw an error")
	assert.Equal(t, true, isUrlYoutube, "absolute valid youtube url was presented")

	isUrlYoutube, err = isYoutubeUrl("https://youtube.com/watch?v=q-cE4ziUBMo")
	assert.Nil(t, err, "valid youtube url should not throw an error")
	assert.Equal(t, true, isUrlYoutube, "absolute valid youtube url was presented")

	isUrlYoutube, err = isYoutubeUrl("https://www.youtube.com/watch")
	assert.Nil(t, err, "valid youtube url without query parameter should not throw an error")
	assert.Equal(t, true, isUrlYoutube, "absolute valid youtube url was presented")
}

func TestGetCmdArgs(t *testing.T) {
	logger := zap.Must(zap.NewProduction())

	os.Args = []string{"main"} // always must be at least one element - programme name
	args, err := getCmdArgs(logger)
	assert.Equal(t, e.ErrIncorrectLinks, err, "need to pass video links via command line arguments")
	assert.Nil(t, args, "if an error occurs it will return uninitialized slice")

	os.Args = []string{"main", "eshelon", "bmstu"}
	args, err = getCmdArgs(logger)
	assert.Equal(t, e.ErrIncorrectLinks, err, "need to pass video links via command line arguments")
	assert.Nil(t, args, "if an error occurs it will return uninitialized slice")

	expected := []string{"https://www.youtube.com/watch"}
	os.Args = []string{"main", "eshelon", "https://www.youtube.com/watch"}
	args, err = getCmdArgs(logger)
	assert.Nil(t, err, "expected that function will return slice of 1 element")
	assert.Equal(t, expected, args)

	expected = []string{"https://www.youtube.com/watch", "https://www.youtube.com/watch?v=q-cE4ziUBMo"}
	os.Args = []string{"main", "https://www.youtube.com/watch", "https://www.youtube.com/watch?v=q-cE4ziUBMo"}
	args, err = getCmdArgs(logger)
	assert.Nil(t, err, "expected that function will return slice of 2 elements")
	assert.Equal(t, expected, args)
}

func TestStartEngine(t *testing.T) {
	ctrl := gomock.NewController(t) // New in go1.14+, if you are passing a *testing.T into this function you no longer need to call ctrl.Finish() in your test methods.
	downloadManagerClientMock := mocks.NewMockDownloadManagerClient(ctrl)
	logger := zap.Must(zap.NewProduction())

	os.Args = []string{"main"}
	err := StartEngine(downloadManagerClientMock, logger)
	assert.Equal(t, e.ErrIncorrectLinks, err, "need to pass video links via command line arguments")

	os.Args = []string{"main", "eshelon"}
	err = StartEngine(downloadManagerClientMock, logger)
	assert.Equal(t, e.ErrIncorrectLinks, err, "need to pass video links via command line arguments")

	os.Args = []string{"main", "eshelon", "https://www.youtube.com/watch"}
	downloadManagerClientMock.EXPECT().Download(gomock.Any(), gomock.Any()).Return(nil, nil)
	err = StartEngine(downloadManagerClientMock, logger)
	assert.Nil(t, err, "expected that function will return nil")
}
