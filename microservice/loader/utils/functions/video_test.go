//go:generate go test ./... -coverprofile=coverage_video.out
//go:generate go tool cover -html=coverage_video.out -o coverage_video.html
package functions

import (
	"testing"
	"unicode/utf8"

	e "github.com/cantylv/thumbnail-loader/microservice/loader/utils/myerrors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFindSubstringIndex(t *testing.T) {
	index := findSubstringIndex("privet", "")
	assert.Equal(t, -1, index)
	index = findSubstringIndex("", "psd")
	assert.Equal(t, -1, index)
	index = findSubstringIndex("privet", "pr")
	assert.Equal(t, 0, index)
	index = findSubstringIndex("privet", "privet")
	assert.Equal(t, 0, index)
	index = findSubstringIndex("https://www.youtube.com/watch?v=q-cE4ziUBMo", "v=")
	assert.Equal(t, utf8.RuneCount([]byte("https://www.youtube.com/watch?v"))-1, index)
}

func TestGetQueryParameter(t *testing.T) {
	videoId := getQueryParameter("", "v")
	assert.Equal(t, "", videoId)
	videoId = getQueryParameter("https://www.youtube.com/watch?", "v")
	assert.Equal(t, "", videoId)
	videoId = getQueryParameter("https://www.youtube.com/watch?v=q-cE4ziUBMo", "v")
	assert.Equal(t, "q-cE4ziUBMo", videoId)
	videoId = getQueryParameter("https://www.youtube.com/watch?v=q-cE4ziUBMo&", "v")
	assert.Equal(t, "q-cE4ziUBMo", videoId)
}

func TestGetVideosIdIncorrectArgs(t *testing.T) {
	parseUris := []string{"https://www.youtube.com/watch"}
	res, err := GetVideosId(parseUris, zap.Must(zap.NewProduction()))
	assert.Nil(t, res, "if error occur, result array will be nil")
	assert.Equal(t, e.ErrIncorrectLinks, err, "incorrect urls were passed")
}

func TestGetVideosIdCorrectArgs(t *testing.T) {
	parseUris := []string{"https://www.youtube.com/watch", "https://www.youtube.com/watch?v=q-cE4ziUBMo"}
	res, err := GetVideosId(parseUris, zap.Must(zap.NewProduction()))
	assert.Nil(t, err, "there must me no error cos of good input")
	expected := []string{"q-cE4ziUBMo"}
	assert.Equal(t, expected, res, "incorrect urls were passed")
}
