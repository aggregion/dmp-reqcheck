package utils

import "net/url"

// MustURLParse helps to parse URL.
func MustURLParse(urlStr string) *url.URL {
	if len(urlStr) == 0 {
		return nil
	}

	var err error
	var result *url.URL

	result, err = url.Parse(urlStr)
	if err != nil {
		panic(err)
	}

	return result
}
