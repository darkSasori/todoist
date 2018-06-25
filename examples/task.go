package main

import (
	"fmt"
	"os"

	"github.com/darkSasori/todoist"
)

func main() {
	todoist.Token = os.Getenv("todoist_token")

	fmt.Println("Create Project")
	project, err := todoist.CreateProject("test")
	if err != nil {
		panic(err)
	}
	fmt.Println(project)

	fmt.Println("List Tasks")
	qp := todoist.QueryParam{}
	tasks, err := todoist.ListTask(qp)
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

	fmt.Println("Create Task")
	task, err := todoist.CreateTask(todoist.Task{Content: "test", ProjectID: project.ID})
	if err != nil {
		panic(err)
	}
	fmt.Println(task)

	fmt.Println("List Tasks project")
	qp["project_id"] = fmt.Sprintf("%d", project.ID)
	tasks, err = todoist.ListTask(qp)
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

	fmt.Println("Update Task")
	task.Content = task.Content + " Update"
	if err := task.Update(); err != nil {
		panic(err)
	}

	fmt.Println("Get Task")
	task, err = todoist.GetTask(task.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(task)

	fmt.Println("Close Task")
	if err := task.Close(); err != nil {
		panic(err)
	}
	tasks, err = todoist.ListTask(qp)
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

	fmt.Println("Reopen Task")
	if err := task.Reopen(); err != nil {
		panic(err)
	}
	tasks, err = todoist.ListTask(qp)
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)

	fmt.Println("Delete Task")
	if err := task.Delete(); err != nil {
		panic(err)
	}

	fmt.Println("Delete Project")
	if err := project.Delete(); err != nil {
		panic(err)
	}

	fmt.Println("End Task")
}
