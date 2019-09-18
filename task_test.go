package todoist

import (
	"encoding/json"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestQueryParam(t *testing.T) {
	qp := QueryParam{}
	if qp.String() != "" {
		t.Errorf("Expected '' != '%s'", qp)
	}

	qp = QueryParam{"param1": "param1"}
	if qp.String() != "?param1=param1" {
		t.Errorf("Expected '?param1=param1' != '%s'", qp)
	}

	qp = QueryParam{
		"param1": "param1",
		"param2": "param2",
	}
	if qp.String() != "?param1=param1&param2=param2" {
		t.Errorf("Expected '?param1=param1&param2=param2' != '%s'", qp)
	}
}

func TestTask(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	listTask := make([]Task, 0)
	var tempTask Task

	httpmock.RegisterResponder("POST", todoistURL+"tasks",
		func(req *http.Request) (*http.Response, error) {
			var body taskSave
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			task := Task{
				ID:        1,
				Content:   body.Content,
				ProjectID: body.ProjectID,
				Order:     body.Order,
				LabelIDs:  body.LabelIDs,
				Priority:  body.Priority,
				Due: Due{
					String:   body.DueString,
					Datetime: body.DueDateTime,
				},
			}
			listTask = append(listTask, task)
			return httpmock.NewJsonResponse(201, task)
		})
	httpmock.RegisterResponder("GET", todoistURL+"tasks",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, listTask)
		})
	httpmock.RegisterResponder("GET", todoistURL+"tasks/1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, listTask[0])
		})
	httpmock.RegisterResponder("POST", todoistURL+"tasks/1",
		func(req *http.Request) (*http.Response, error) {
			var body taskSave
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			task := Task{
				ID:        1,
				Content:   body.Content,
				ProjectID: body.ProjectID,
				Order:     body.Order,
				LabelIDs:  body.LabelIDs,
				Priority:  body.Priority,
				Due: Due{
					String:   body.DueString,
					Datetime: body.DueDateTime,
				},
			}
			listTask[0] = task
			return httpmock.NewJsonResponse(200, task)
		})
	httpmock.RegisterResponder("POST", todoistURL+"tasks/1/close",
		func(req *http.Request) (*http.Response, error) {
			tempTask = listTask[0]
			listTask = make([]Task, 0)
			return httpmock.NewJsonResponse(200, nil)
		})
	httpmock.RegisterResponder("POST", todoistURL+"tasks/1/reopen",
		func(req *http.Request) (*http.Response, error) {
			listTask = append(listTask, tempTask)
			return httpmock.NewJsonResponse(200, nil)
		})
	httpmock.RegisterResponder("DELETE", todoistURL+"tasks/1",
		func(req *http.Request) (*http.Response, error) {
			tempTask = listTask[0]
			listTask = make([]Task, 0)
			return httpmock.NewJsonResponse(200, nil)
		})

	task, err := CreateTask(Task{Content: "test"})
	if err != nil {
		t.Error(err)
	}
	if task.Content != "test" {
		t.Errorf("Expected 'test' != '%s'", task.Content)
	}

	tasks, err := ListTask(QueryParam{})
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected len(tasks) != '%d'", len(tasks))
	}

	task.Content = "updated"
	if err := task.Update(); err != nil {
		t.Error(err)
	}

	task, err = GetTask(1)
	if err != nil {
		t.Error(err)
	}
	if task.Content != "updated" {
		t.Errorf("Expected 'updated' != '%s'", task.Content)
	}

	if err := task.Close(); err != nil {
		t.Error(err)
	}

	tasks, err = ListTask(QueryParam{})
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected len(tasks) != '%d'", len(tasks))
	}

	if err := task.Reopen(); err != nil {
		t.Error(err)
	}

	tasks, err = ListTask(QueryParam{})
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected len(tasks) != '%d'", len(tasks))
	}

	if err := task.Delete(); err != nil {
		t.Error(err)
	}
}
