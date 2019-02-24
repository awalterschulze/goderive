package toerror

import (
	"fmt"
	"net/http"
)

func parseMinorHTTPVersion(versionString string) (int, error) {
	return deriveCompose(
		deriveToError(fmt.Errorf("HTTP version parsing failed"), http.ParseHTTPVersion),
		func(major, minor int) (int, error) {
			if major != 2 {
				return 0, fmt.Errorf("only HTTP2 is supported")
			}
			return minor, nil
		},
	)(versionString)
}
