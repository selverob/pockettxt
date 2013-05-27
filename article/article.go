package article

import (
	"bytes"
)

type Article struct {
	Title  string
	Author string
	Text   string
	//Date could be time.Time but
	//since we're not doing anything with it,
	//why bother converting it?
	Date string
	URL  string
}

func (a Article) Print() (o *bytes.Buffer) {
	o = &bytes.Buffer{}
	o.WriteString(a.Title)
	o.WriteString("\nby ")
	o.WriteString(a.Author)
	o.WriteString("\n\n")
	o.WriteString(a.Text)
	o.WriteString("\n\nPublished on: ")
	o.WriteString(a.Date)
	o.WriteString("\nFound at: ")
	o.WriteString(a.URL)
	return
}

func PrintArticles(as []Article) (o *bytes.Buffer) {
	o = &bytes.Buffer{}
	for _, a := range as {
		a.Print().WriteTo(o)
		o.WriteString("\n\n=====================================\n\n")
	}
	return
}
