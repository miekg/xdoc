package main

import (
	"fmt"
	"log"
	"net/url"

	gu "github.com/miekg/xdoc/gitlabutil"
)

const DocDir = "docs"

func main() {
	base := "https://gitlab.gnome.org"
	cl := gu.NewClient(base)
	groups, err := gu.ListGroups(cl)
	if err != nil {
		log.Fatal(err)
	}

	if len(groups) == 0 {
		log.Fatal("No groups found")
	}

	gid := groups[0].ID
	proj, err := gu.ListProjects(cl, gid)
	if err != nil {
		log.Fatal(err)
	}
	println(proj[0].WebURL)
	url, _ := url.Parse(proj[0].WebURL)
	println(url.Path)
	println(len(proj), "found")
	files, _ := gu.ListDir(cl, proj[0].ID, DocDir)
	fmt.Printf("%+v\n", files)
}
