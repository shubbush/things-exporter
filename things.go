package main

import (
	"database/sql"
	"log"
)

type ThingsDB struct {
	db *sql.DB
}

type Tag struct {
	uuid  string
	title string
	index int
}

type Area struct {
	uuid  string
	title string
	index int
}

type Task struct {
	uuid        string
	trashed     int
	taskType    int // 0 - task, 1 - project, 2 - heading
	title       string
	notes       *string
	dueDate     *int64
	status      int
	startDate   *int64
	index       int
	todayIndex  int
	area        string
	project     string
	actionGroup string // heading
}

type CheckListItem struct {
	uuid   string
	title  string
	status int // todo: check corresponding values
	index  int
	task   string
}

func (t ThingsDB) getTags() ([]Tag, error) {
	rows, err := t.db.Query("select uuid, title, \"index\" from TMTag")
	if err != nil {
		log.Fatal("Error during fetching tags: ", err)
		return nil, err
	}
	defer rows.Close()

	data := []Tag{}
	for rows.Next() {
		t := Tag{}
		err = rows.Scan(&t.uuid, &t.title, &t.index)
		if err != nil {
			return nil, err
		}
		data = append(data, t)
	}
	return data, nil
}
