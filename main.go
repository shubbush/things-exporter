package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niklasfasching/go-org/org"
)

const dbFile string = "db/main.sqlite"

const outDir = "out"

var emptyTags = make([]string, 0)

func main() {

	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	thingsDb := ThingsDB{db}
	tagsByTask := getTagsByTask(&thingsDb)
	areasByUuid := getAreasMap(&thingsDb)

	projects, err := thingsDb.getProjects()
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		nodes := []org.Node{}
		// fmt.Printf("Processing %s\n", project.Title)
		headline := convertProject(&thingsDb, project, tagsByTask)
		nodes = append(nodes, headline)

		doc := org.Document{
			Nodes: nodes,
		}
		orgWriter := org.OrgWriter{}
		docStr, err := doc.Write(&orgWriter)
		if err != nil {
			panic(err)
		}
		area := ""
		if project.Area.Valid {
			area = areasByUuid[project.Area.String]
		}
		path := mkDirPath(area)
		os.MkdirAll(path, os.ModePerm)
		fileName := mkFileName(project.Title)
		file, err := os.Create(filepath.Join(path, fileName))
		file.Write([]byte(docStr))
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

}

func mkDirPath(area string) string {
	path := filepath.Join(".", outDir)
	if len(area) > 0 {
		path = filepath.Join(path, area)
	}
	return path
}

func mkFileName(project string) string {
	return fmt.Sprintf("%s.org", strings.ReplaceAll(strings.ToLower(project), "\\s", "_"))
}

func getTagsByTask(db *ThingsDB) map[string][]string {
	tags, err := db.getTagsByTask()
	if err != nil {
		panic(err)
	}
	res := make(map[string][]string)
	for _, tag := range tags {
		res[tag.Uuid] = strings.Split(tag.Tags, ",")
	}
	return res
}

func getAreasMap(db *ThingsDB) map[string]string {
	areas := db.getAreas()
	res := make(map[string]string)
	for _, area := range areas {
		res[area.Uuid] = area.Title
	}
	return res
}

func convertProject(db *ThingsDB, thingsProject Task, tagsByTaskUuid map[string][]string) org.Headline {
	headline := taskToTodo(1, thingsProject, emptyTags)
	tasks, err := db.getTasksByProject(thingsProject.Uuid)
	if err != nil {
		fmt.Println(err)
	}
	for _, task := range tasks {
		tags := tagsByTaskUuid[task.Uuid]
		todo := taskToTodo(2, task, tags)
		addChild(&headline, todo)
	}
	return headline
}

func taskToTodo(lvl int, task Task, tags []string) org.Headline {
	todo := mkTodo(lvl, task.Title, taskStatusToOrgStatus(task.Status), task.Index, task.DueDate, tags)
	if task.StopDate.Valid {
		sec, dec := math.Modf(task.StopDate.Float64)
		closed := time.Unix(int64(sec), int64(dec*(1e9)))
		formatedClosed := fmt.Sprintf("CLOSED: [%s]", closed.Format("2006-01-02 Mon 15:04"))
		addChild(&todo, mkText(formatedClosed+"\n"))
	}
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
