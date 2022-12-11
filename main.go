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
	tags, err := thingsDb.getTags()
	fmt.Println("Tags:")
	for _, el := range tags {
		fmt.Printf("%s\n", el.Title)
	}
	fmt.Println("")

	// areas := thingsDb.getAreas()
	projects, err := thingsDb.getProjectsWithoutArea()

	fmt.Println(err)

	for _, project := range projects {
		fmt.Printf("Processing %s\n", project.Title)
		headline := mkHeadline(1, project.Title)
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
