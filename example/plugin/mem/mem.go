package mem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name string
}

func getUser(version int, id int) (User, error) {
	resp, err := http.Get(fmt.Sprintf("localhost:8080/user/%d/%d", version, id))
	if err != nil {
		return User{}, nil
	}
	u := &User{}
	if err := json.NewDecoder(resp.Body).Decode(u); err != nil {
		return User{}, nil
	}
	return *u, nil
}

var mgetUser = deriveMem(getUser)

func GetUser(version int, id int) (User, error) {
	return mgetUser(version, id)
}
