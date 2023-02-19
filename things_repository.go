package main

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ThingsDB struct {
	db *sqlx.DB
}

type Tag struct {
	Uuid  string
	Title string
	Index int
}

type Area struct {
	Uuid  string
	Title string
	Index int
}

// taskType    int // 0 - task, 1 - project, 2 - heading
// todayIndex  int
// uuid    string
type Task struct {
	Uuid    string
	Title   string
	Notes   sql.NullString
	DueDate sql.NullString `db:"due_date"`
	Index   int
	Status  int // 0 - todo, 2 - canceled, 3 - completed
	Area    sql.NullString
	Project sql.NullString
}

type TaskWithTags struct {
	Uuid string
	Tags string
}

type CheckListItem struct {
	uuid   string
	title  string
	status int // todo: check corresponding values
	index  int
	task   string
}

func (t ThingsDB) getTags() ([]Tag, error) {
	tags := []Tag{}
	e := t.db.Select(&tags, "select uuid, title, \"index\" from TMTag")
	return tags, e
}

func (t ThingsDB) getTagsByTask() ([]TaskWithTags, error) {
	taskWithTags := []TaskWithTags{}
	e := t.db.Select(&taskWithTags, "SELECT tt.tasks as uuid, group_concat(t.title) as tags from TMTaskTag tt join TMTag t on tt.tags = t.uuid group by tt.tasks")
	return taskWithTags, e
}

func (t ThingsDB) getAreas() []Area {
	areas := []Area{}
	t.db.Select(&areas, "select uuid, title, \"index\" from TMArea")
	return areas
}

func (t ThingsDB) getProjects() ([]Task, error) {
	projects := []Task{}
	err := t.db.Select(&projects, "SELECT t.uuid, t.title, t.notes, t.\"index\", t.status, t.area, t.project from TMTask t where t.\"type\" == 1 and t.trashed != 1 order by t.\"index\"")
	return projects, err
}

func (t ThingsDB) getTasksByProject(projectId string) ([]Task, error) {
	tasks := []Task{}
	err := t.db.Select(&tasks, "SELECT t.uuid, t.title, t.notes, date(t.dueDate, 'unixepoch') AS due_date, t.\"index\", t.status, t.area, t.project from TMTask t where t.\"type\" == 0 and t.trashed != 1 and t.project == ? order by t.\"index\"", projectId)
	return tasks, err
}
