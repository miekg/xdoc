package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	gu "github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
)

const DocDir = "docs"

var (
	flgGroup = flag.String("group", "", "select only this group")
	flgBase  = flag.String("base", "https://gitlab.com", "GitLab site")
	flgDir   = flag.String("dir", DocDir, "directory to use for documentation")
)

func main() {
	flag.Parse()
	cl := gu.NewClient(*flgBase)

	groups, err := gu.ListGroups(cl)
	if err != nil {
		log.Fatal(err)
	}
	if len(groups) == 0 {
		log.Fatal("No groups found")
	}

	names := make([]string, len(groups))
	for i := range groups {
		names[i] = groups[i].Name
	}
	log.Printf("Groups found: %q\n", names)

	var group *gitlab.Group
	for i := range groups {
		if groups[i].Name == *flgGroup {
			group = groups[i]
			break
		}
	}
	if group == nil {
		log.Fatalf("Group %q not found", *flgGroup)
	}

	gid := group.ID
	proj, err := gu.ListProjects(cl, gid)
	if err != nil {
		log.Fatal(err)
	}
	println(proj[0].WebURL)
	url, _ := url.Parse(proj[0].WebURL)
	println(url.Path)
	println(len(proj), "found")
	files, _ := gu.ListDir(cl, proj[0].ID, *flgDir)
	log.Printf("%d files found in %q", len(files), url.Path)
	fmt.Printf("%+v\n", files)
	for i := range files {

	}
}
