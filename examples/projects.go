package main

import (
	"fmt"
	"os"

	"github.com/wtfutil/todoist"
)

func main() {
	todoist.Token = os.Getenv("todoist_token")

	fmt.Println("CreateProject")
	project, err := todoist.CreateProject("teste")
	if err != nil {
		panic(err)
	}
	fmt.Println(project)

	fmt.Println("ListProject")
	projects, err := todoist.ListProject()
	if err != nil {
		panic(err)
	}
	fmt.Println(projects)

	fmt.Println("UpdateProject")
	project.Name = project.Name + " Update"
	if err = project.Update(); err != nil {
		panic(err)
	}

	fmt.Println("GetProject")
	project, err = todoist.GetProject(project.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(project)

	fmt.Println("DeleteProject")
	if err := project.Delete(); err != nil {
		panic(err)
	}

	fmt.Println("End Project")
}
