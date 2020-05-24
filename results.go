package main

import (
	"fmt"

	"github.com/blevesearch/bleve"
)

func (d *Doc) HTML(sr *bleve.SearchResult) string {
	rv := `<html>
<head>
<title>Results</title>
<link rel="stylesheet" type="text/css" href="/a/style.css"/>
</head>

<body>
`
	if sr.Total > 0 {
		if sr.Request.Size > 0 {
			rv = fmt.Sprintf("<p>%d matches, showing %d through %d, took %s</p>", sr.Total, sr.Request.From+1, sr.Request.From+len(sr.Hits), sr.Took)
			rv += "<ol>\n"
			for _, hit := range sr.Hits {
				// hit.ID  is the url
				rv += fmt.Sprintf(`<li><a href="/r/%s">%s</a> (%f)`, hit.ID, hit.ID, hit.Score)
				for fragmentField, fragments := range hit.Fragments {
					println("fragment") // TODO: what is this, when do we see it?
					rv += fmt.Sprintf("\t%s\n", fragmentField)
					for _, fragment := range fragments {
						rv += fmt.Sprintf("\t\t%s\n", fragment)
					}
				}
				for otherFieldName, otherFieldValue := range hit.Fields {
					println("otherfields")
					if _, ok := hit.Fragments[otherFieldName]; !ok {
						rv += fmt.Sprintf("\t%s\n", otherFieldName)
						rv += fmt.Sprintf("\t\t%v\n", otherFieldValue)
					}
				}
				rv += "</li>"
			}
			rv += "</ol>\n"
		} else {
			rv = fmt.Sprintf("%d matches, took %s\n", sr.Total, sr.Took)
		}
	} else {
		rv = "No matches"
	}
	// TODO(miek): wth is this actually?
	if len(sr.Facets) > 0 {
		rv += fmt.Sprintf("Facets:\n")
		for fn, f := range sr.Facets {
			rv += fmt.Sprintf("%s(%d)\n", fn, f.Total)
			for _, t := range f.Terms {
				rv += fmt.Sprintf("\t%s(%d)\n", t.Term, t.Count)
			}
			if f.Other != 0 {
				rv += fmt.Sprintf("\tOther(%d)\n", f.Other)
			}
		}
	}
	return rv + `</body></html>`
}
