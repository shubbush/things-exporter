package main

import (
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
	Title   string
	Notes   *string
	Index   int
	Status  int // 0 - todo, 2 - canceled, 3 - completed
	Area    *string
	Project *string
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

func (t ThingsDB) getAreas() []Area {
	areas := []Area{}
	t.db.Select(&areas, "select uuid, title, \"index\" from TMArea")
	return areas
}

func (t ThingsDB) getProjectsWithoutArea() ([]Task, error) {
	projects := []Task{}
	err := t.db.Select(&projects, "SELECT t.title, t.notes, t.\"index\", t.status, t.area, t.project from TMTask t where t.\"type\" == 1 and t.trashed != 1 and t.area is null order by t.\"index\"")
	return projects, err
}
