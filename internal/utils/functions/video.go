package functions

import (
	"fmt"
	"os"
	"regexp"

	e "github.com/cantylv/thumbnail-loader/internal/utils/myerrors"
	"go.uber.org/zap"
)

var (
	hasDomain = `^[a-z/:]*[www.]?youtube.com+`
)

// GetVideosId
// returns the value of query parameter 'v' from urls
func getVideosId(logger *zap.Logger) ([]string, error) {
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

	videosId := make([]string, 0, len(parseUris))
	for _, uri := range parseUris {
		videoID := getQueryParameter(uri, "v")
		if videoID != "" {
			videosId = append(videosId, videoID)
		}
	}
	if len(videosId) == 0 {
		return nil, e.ErrIncorrectLinks
	}
	return videosId, nil
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

// getQueryParameter
// returns the value of the passed query parameter
func getQueryParameter(url, parameterName string) string {
	startIndex := findSubstringIndex(url, parameterName+"=")
	if startIndex == -1 {
		return ""
	}
	startIndex += 2 // +2, because we want to skip '='
	var endIndex = startIndex
	for ; endIndex < len(url); endIndex++ {
		if url[endIndex] == '&' {
			break
		}
	}
	return url[startIndex:endIndex]
}

// findSubstringIndex
// returns the index of query parametr (index of first symbol of query parameter name)
func findSubstringIndex(str, subStr string) (startIndex int) {
	var i, j int
	var isCheckingMatch bool
	for ; i < len(str) && j < len(subStr); i++ {
		if str[i] == subStr[j] {
			if !isCheckingMatch {
				startIndex = i
				isCheckingMatch = true
			}
			j++
			continue
		}
		isCheckingMatch = false
		j = 0
	}
	if i == len(str) && j != len(subStr) {
		return -1
	}
	return startIndex
}
