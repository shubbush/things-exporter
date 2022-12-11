package main

import "github.com/niklasfasching/go-org/org"

const TODO = "TODO"
const DONE = "DONE"

func addChild(parent *org.Headline, child org.Node) {
	parent.Children = append(parent.Children, child)
}

func mkHeadline(lvl int, title string) org.Headline {
	return org.Headline{Lvl: lvl, Title: []org.Node{org.Text{Content: title, IsRaw: true}}, Children: []org.Node{}}
}

func mkTodo(lvl int, title string, status string) org.Headline {
	return org.Headline{
		Lvl:    lvl,
		Title:  []org.Node{org.Text{Content: title, IsRaw: true}},
		Status: status}
}
