package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", todoistURL+"test",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, nil)
		})
	httpmock.RegisterResponder("DELETE", todoistURL+"test",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, "not found")
		})
	httpmock.RegisterResponder("POST", todoistURL+"test",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewJsonResponse(200, body)
		})

	res, err := makeRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected Status 200 != '%d'", res.StatusCode)
	}

	res, err = makeRequest("POST", "test", struct{ ID int }{1})
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected Status 200 != '%d'", res.StatusCode)
	}
	var body struct{ ID int }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Error(err)
	}
	if body.ID != 1 {
		t.Errorf("Expected 1 != '%d'", body.ID)
	}

	res, err = makeRequest("DELETE", "test", nil)
	if err.Error() != "\"not found\"" {
		t.Error(err)
	}
}

func TestCustomTime(t *testing.T) {
	var time CustomTime
	b, _ := time.MarshalJSON()
	if string(b) != "null" {
		t.Errorf("Expected 'null' != '%s'", string(b))
	}

	time.UnmarshalJSON([]byte("null"))
	if !time.IsZero() {
		t.Errorf("Expected time is zero")
	}

	if err := time.UnmarshalJSON([]byte("2018-07-01T01:00:00+00:00")); err != nil {
		t.Error(err)
	}

	b, err := time.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "\"2018-07-01T01:00:00+00:00\"" {
		t.Errorf("Expected '2018-07-01T01:00:00+00:00' != '%s'", string(b))
	}
}

func TestTaskSaveMarshalJSON(t *testing.T) {
	var ts taskSave
	_, err := ts.MarshalJSON()
	if err.Error() != "Content is empty" {
		t.Error(err)
	}

	ts.Content = "test"
	b, _ := ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\"}" {
		t.Errorf("Expected '{\"content\":\"test\"}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", ProjectID: 1}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"project_id\":1}" {
		t.Errorf("Expected '{\"content\":\"test\",\"project_id\":1}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", Order: 1}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"order\":1}" {
		t.Errorf("Expected '{\"content\":\"test\",\"order\":1}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", Priority: 1}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"priority\":1}" {
		t.Errorf("Expected '{\"content\":\"test\",\"priority\":1}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", DueString: "en"}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"due_string\":\"en\"}" {
		t.Errorf("Expected '{\"content\":\"test\",\"due_string\":\"en\"}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", DueLang: "en"}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"due_lang\":\"en\"}" {
		t.Errorf("Expected '{\"content\":\"test\",\"due_lang\":\"en\"}' != '%s'", string(b))
	}

	ts = taskSave{Content: "test", LabelIDs: []int{1, 2}}
	b, _ = ts.MarshalJSON()
	if string(b) != "{\"content\":\"test\",\"label_ids\":[1,2]}" {
		t.Errorf("Expected '{\"content\":\"test\",\"label_ids\":[1,2]}' != '%s'", string(b))
	}

	now := CustomTime{time.Now()}
	ts = taskSave{Content: "test", DueDateTime: now}
	s, _ := now.MarshalJSON()
	b, _ = ts.MarshalJSON()
	if string(b) != fmt.Sprintf("{\"content\":\"test\",\"due_datetime\":%s}", string(s)) {
		t.Error("deu merda")
	}
}
