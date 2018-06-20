package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

const ctLayout = "2006-01-02T15:04:05+00:00"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}

	ct.Time, err = time.Parse(ctLayout, s)
	return err
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format(ctLayout) + `"`), nil
}
