package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xeonx/timeago"
)

type Story struct {
	Url       string
	Title     string
	Guid      string
	Ago       string
	Timestamp string
	Username  string
	HasUrl    bool
	Domain    string
}

type WelcomeVars struct {
	Rows []*Story
}

const layout = "2006-01-02 15:04:05"

func WelcomeIndexVars(db *sqlx.DB, location *time.Location) *WelcomeVars {
	vars := WelcomeVars{}
	vars.Rows = []*Story{}
	rows, _ := db.Queryx("SELECT * FROM stories ORDER BY created_at desc limit 30")
	for rows.Next() {
		m := make(map[string]any)
		rows.MapScan(m)
		story := storyFromMap(m)
		vars.Rows = append(vars.Rows, story)
	}
	return &vars
}

func storyFromMap(m map[string]any) *Story {
	story := Story{}
	story.Title = fmt.Sprintf("%s", m["title"])
	story.Url = fmt.Sprintf("%s", m["url"])
	story.Guid = fmt.Sprintf("%s", m["guid"])
	if story.Url != "" {
		story.HasUrl = true
		tokens := strings.Split(story.Url, "/")
		if len(tokens) > 2 {
			tokens = strings.Split(tokens[2], ".")
			if len(tokens) == 3 {
				tokens = tokens[1:]
			}
			story.Domain = strings.Join(tokens, ".")
		}
	}

	tm := m["created_at"].(time.Time)
	tm = tm.Add(time.Hour * 8)
	story.Timestamp = fmt.Sprintf("%s", tm)
	story.Ago = timeago.English.Format(tm)
	return &story
}
