package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/miekg/xdoc/gitlabutil"
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

	doc := New()
	doc.Insert(proj[0])
	files, _ := gu.ListDir(cl, proj[0].ID, *flgDir)
	fmt.Printf("%d\n", len(files))
	for i := range files {
		log.Printf("Downloading %q %s", path.Join(proj[0].WebURL, files[i].Path), files[i].Type)
		buf, err := gitlabutil.Download(cl, proj[0].ID, "master", files[i].Path)
		if err != nil {
			log.Fatal(err)
		}
		doc.InsertFile(proj[0], files[i].Path, buf)
	}

	r := doc.setup()
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
