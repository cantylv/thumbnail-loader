package myerrors

import "errors"

var (
	ErrIncorrectLinks = errors.New("you've provided incorrect links to the video, please try again; for correct utilite work you need to pass URIs of videos separated by space; URIs must have query parameter \"v\"")
)
