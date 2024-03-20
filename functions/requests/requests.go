package requests

import (
	"encoding/json"
	"io"
	"net/http"
)

func Get(token string, url string, result interface{}) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// headers
	req.Header.Set("Authorization", "Bearer "+token)

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// read request
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// parse json
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}
