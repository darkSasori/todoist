package todoist

import (
	"encoding/json"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestProject(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	listProjects := make([]Project, 0)

	httpmock.RegisterResponder("POST", todoistURL+"projects",
		func(req *http.Request) (*http.Response, error) {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			project := Project{
				ID:   1,
				Name: body.Name,
			}
			listProjects = append(listProjects, project)
			return httpmock.NewJsonResponse(201, project)
		})
	httpmock.RegisterResponder("GET", todoistURL+"projects",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, listProjects)
		})
	httpmock.RegisterResponder("GET", todoistURL+"projects/1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, listProjects[0])
		})
	httpmock.RegisterResponder("POST", todoistURL+"projects/1",
		func(req *http.Request) (*http.Response, error) {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			listProjects[0].Name = body.Name
			return httpmock.NewJsonResponse(200, listProjects[0])
		})
	httpmock.RegisterResponder("DELETE", todoistURL+"projects/1",
		func(req *http.Request) (*http.Response, error) {
			listProjects = make([]Project, 0)
			return httpmock.NewJsonResponse(200, "")
		})

	project, err := CreateProject("test")
	if err != nil {
		t.Error(err)
	}
	if project.Name != "test" {
		t.Errorf("Expected 'test' != '%s'", project.Name)
	}

	projects, err := ListProject()
	if err != nil {
		t.Error(err)
	}
	if len(projects) != 1 {
		t.Errorf("Expected len(projects) != '%d'", len(projects))
	}

	project, err = GetProject(1)
	if err != nil {
		t.Error(err)
	}
	if project.Name != "test" {
		t.Errorf("Expected 'test' != '%s'", project.Name)
	}

	project.Name = "updated"
	err = project.Update()
	if err != nil {
		t.Error(err)
	}
	if project.Name != "updated" {
		t.Errorf("Expected 'updated' != '%s'", project.Name)
	}

	err = project.Delete()
	if err != nil {
		t.Error(err)
	}
}
