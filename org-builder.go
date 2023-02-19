package main

import (
	"database/sql"

	"github.com/niklasfasching/go-org/org"
)

const TODO = "TODO"
const DONE = "DONE"
const CANCELED = "CANCELED"

func addChild(parent *org.Headline, child org.Node) {
	parent.Children = append(parent.Children, child)
}

func mkHeadline(lvl int, title string) org.Headline {
	return org.Headline{Lvl: lvl, Title: []org.Node{org.Text{Content: title, IsRaw: true}}, Children: []org.Node{}}
}

func mkTodo(lvl int, title string, status string, index int, dueDate sql.NullString, tags []string) org.Headline {
	return org.Headline{
		Index:  index,
		Lvl:    lvl,
		Title:  []org.Node{mkText(title)},
		Status: status,
		Tags:   tags,
	}
}

func mkText(text string) org.Text {
	return org.Text{Content: text, IsRaw: false}
}
