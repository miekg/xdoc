package main

import (
	"flag"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/miekg/xdoc/gitlabutil"
	gu "github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
)

const DocDir = "docs"

var (
	flgGroup = flag.String("group", "ALL", "select only this group")
	flgBase  = flag.String("base", "https://gitlab.com", "GitLab site")
	flgDir   = flag.String("dir", DocDir, "directory to use for documentation")
)

func main() {
	flag.Parse()
	cl := gu.NewClient(*flgBase)

	groups := []*gitlab.Group{}
	var err error
	if *flgGroup == "ALL" {
		if groups, err = gu.ListGroups(cl); err != nil {
			log.Fatal(err)
		}
	} else {
		if groups, err = gu.GetGroup(cl, *flgGroup); err != nil {
			log.Fatal(err)
		}
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

	// TODO(miek): wrap the stuff below in a loop and do this continuously
	doc := New()
	doc.Insert(proj[0])
	files, _ := gu.ListDir(cl, proj[0].ID, *flgDir)
	for i := range files {
		log.Printf("Downloading %q %s\n", path.Join(proj[0].WebURL, files[i].Path), files[i].Type)
		buf, err := gitlabutil.Download(cl, proj[0].ID, "master", files[i].Path)
		if err != nil {
			log.Fatal(err)
		}
		doc.InsertFile(proj[0], files[i].Path, buf)
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Created index, now indexing")
	for _, g := range doc.Projects {
		for k, buf := range g.Files {
			if err := index.Index(k, string(buf)); err != nil {
				log.Fatal(err)
			}
		}

	}
	doc.Index = index
	count, err := index.DocCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Downloaded and indexed %d files, starting web server\n", count)

	r := doc.setup()
	println(doc.String())
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
