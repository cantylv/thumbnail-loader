package functions

import (
	"unicode/utf8"

	e "github.com/cantylv/thumbnail-loader/microservice/loader/utils/myerrors"
	"go.uber.org/zap"
)

// GetVideosId
// returns the value of query parameter 'v' from urls
func GetVideosId(parseUris []string, logger *zap.Logger) ([]string, error) {
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
	if utf8.RuneCount([]byte(subStr)) == 0 || utf8.RuneCount([]byte(str)) == 0 {
		return -1
	}
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
