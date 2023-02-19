package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niklasfasching/go-org/org"
)

const dbFile string = "db/main.sqlite"

var emptyTags = make([]string, 0)

func main() {

	nodes := []org.Node{}

	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	thingsDb := ThingsDB{db}
	tagsByTask := getTagsByTask(&thingsDb)

	projects, err := thingsDb.getProjects()
	if err != nil {
		fmt.Println(err)
	}

	for _, project := range projects {
		// fmt.Printf("Processing %s\n", project.Title)
		headline := convertProject(&thingsDb, project, tagsByTask)
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

func getTagsByTask(db *ThingsDB) map[string][]string {
	tags, err := db.getTagsByTask()
	if err != nil {
		log.Fatal(err)
	}
	res := make(map[string][]string)
	for _, tag := range tags {
		res[tag.Uuid] = strings.Split(tag.Tags, ",")
	}
	return res
}

func convertProject(db *ThingsDB, thingsProject Task, tagsByTaskUuid map[string][]string) org.Headline {
	headline := taskToTodo(0, thingsProject, emptyTags)
	tasks, err := db.getTasksByProject(thingsProject.Uuid)
	if err != nil {
		fmt.Println(err)
	}
	for _, task := range tasks {
		tags := tagsByTaskUuid[task.Uuid]
		todo := taskToTodo(1, task, tags)
		addChild(&headline, todo)
	}
	return headline
}

func taskToTodo(lvl int, task Task, tags []string) org.Headline {
	todo := mkTodo(lvl, task.Title, taskStatusToOrgStatus(task.Status), task.Index, task.DueDate, tags)
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
		orgStatus = DONE // todo: support timestamp?
	}
	return orgStatus
}
