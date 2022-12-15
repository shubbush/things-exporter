package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niklasfasching/go-org/org"
)

const dbFile string = "db/main.sqlite"

func main() {

	nodes := []org.Node{}

	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	thingsDb := ThingsDB{db}
	// tags, err := thingsDb.getTags()

	// areas := thingsDb.getAreas()
	projects, err := thingsDb.getProjectsWithoutArea()

	if err != nil {
		fmt.Println(err)
	}

	for _, project := range projects {
		// fmt.Printf("Processing %s\n", project.Title)
		headline := taskToTodo(1, project)
		tasks, err := thingsDb.getTasksByProject(project.Uuid)
		if err != nil {
			fmt.Println(err)
		}
		for _, task := range tasks {
			todo := taskToTodo(2, task)
			addChild(&headline, todo)
		}
		nodes = append(nodes, headline)
	}

	doc := org.Document{
		Nodes: nodes,
	}
	orgWriter := org.OrgWriter{}
	docStr, err := doc.Write(&orgWriter)
	fmt.Println(docStr)

	if err != nil {
		log.Fatal(err)
	}

}

func taskToTodo(lvl int, task Task) org.Headline {
	todo := mkTodo(lvl, task.Title, taskStatusToOrgStatus(task.Status), task.Index, task.DueDate)
	if task.Notes.Valid && len(task.Notes.String) > 0 {
		addChild(&todo, mkText(task.Notes.String+"\n"))
	}
	return todo
}

func taskStatusToOrgStatus(status int) string {
	orgStatus := TODO
	if status == 2 {
		orgStatus = CANCELED
	} else if status == 3 {
		orgStatus = DONE
	}
	return orgStatus
}
