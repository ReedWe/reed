package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Call(path string, request interface{}) error {

	var bodyReader io.Reader
	if request != nil {
		var jsonBody bytes.Buffer
		if err := json.NewEncoder(&jsonBody).Encode(request); err != nil {
			fmt.Printf("err %v \n", err)
		}
		bodyReader = &jsonBody
	}

	u, err := url.Parse("http://127.0.0.1:9888")
	if err != nil {
		fmt.Printf("%v \n", err)
	}

	u.Path = path

	req, err := http.NewRequest("POST", u.String(), bodyReader)
	if err != nil {
		fmt.Printf("%v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.DefaultClient
	r, err := (*client).Do(req)
	if err != nil {
		fmt.Printf("%v", err)
	}

	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return errors.New("rpc error")
	}

	return nil

}
