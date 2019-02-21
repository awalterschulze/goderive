package toerror

import (
	"fmt"
	"net/http"
)

func parseMajorHTTPVersion(versionString string) (int, error) {
	return deriveCompose(
		deriveToError(fmt.Errorf("HTTP version parsing failed"), http.ParseHTTPVersion),
		func(major, minor int) (int, error) {
			return major, nil
		},
	)(versionString)
}
