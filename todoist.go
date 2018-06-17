package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Token save the personal token from todoist
var Token string
var todoistURL = "https://beta.todoist.com/API/v8/"

func makeRequest(method, endpoint string, data interface{}) (*http.Response, error) {
	url := todoistURL + endpoint
	body := bytes.NewBuffer([]byte{})

	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(json)
	}

	fmt.Println(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	bearer := fmt.Sprintf("Bearer %s", Token)
	req.Header.Add("Authorization", bearer)

	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		defer res.Body.Close()
		str, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(str))
	}

	return res, nil
}
