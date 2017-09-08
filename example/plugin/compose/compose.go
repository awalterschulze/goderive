package compose

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func upgradeUser(endpoint, username string) error {
	getEndpoint := fmt.Sprintf("%s/oldusers/%s", endpoint, username)
	postEndpoint := fmt.Sprintf("%s/newusers/%s", endpoint, username)

	_, err := deriveCompose(
		http.Get,
		func(req *http.Response) ([]byte, error) {
			return ioutil.ReadAll(req.Body)
		},
		newUserFromJson,
		newUserFromUser,
		json.Marshal,
		func(buf []byte) (*http.Response, error) {
			return http.Post(
				postEndpoint,
				"application/json",
				bytes.NewBuffer(buf),
			)
		},
	)(getEndpoint)
	return err
}

type user struct {
	Name string
}

func newUserFromJson(buf []byte) (*user, error) {
	u := &user{}
	return u, json.Unmarshal(buf, u)
}

type newUser struct {
	Name      string
	LastNames string
}

func newUserFromUser(u *user) (*newUser, error) {
	names := strings.Split(u.Name, " ")
	n := &newUser{Name: names[0]}
	if len(names) <= 1 {
		return nil, fmt.Errorf("user has no last name")
	}
	n.LastNames = strings.Join(names[1:], " ")
	return n, nil
}
