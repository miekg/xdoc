package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/miekg/xdoc/gitlabutil"
	gu "github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
)

const DocDir = "docs"

var (
	flgGroup = flag.String("group", "", "select only this group")
	flgBase  = flag.String("base", "https://gitlab.com", "GitLab site")
	flgDir   = flag.String("dir", DocDir, "directory to use for documentation")
	flgInt   = flag.Duration("int", 10*time.Minute, "duration to sleep before restarting the download loop")
)

func main() {
	flag.Parse()
	if *flgGroup == "" {
		log.Fatal("-group need a value")
	}

	cl := gu.NewClient(*flgBase)
	// all this group stuff is here to support multiple group, the POC won't do this (yet).
	groups, err := gu.GetGroup(cl, *flgGroup)
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

	doc := New()

	gid := group.ID
	projs, err := gu.ListProjects(cl, gid)
	if err != nil {
		log.Fatal(err)
	}

	if err := doc.InsertProjects(cl, projs); err != nil {
		log.Fatal(err)
	}
	println(doc.String())

	r := doc.setup()

	go func() {
		srv := &http.Server{
			Handler:      r,
			Addr:         "127.0.0.1:8000",
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	tick := time.NewTicker(*flgInt)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-tick.C:
			projs, err := gu.ListProjects(cl, gid)
			if err != nil {
				log.Fatal(err)
			}

			// how do delete deleted projects? Separate list of names:deleted
			if err := doc.InsertProjects(cl, projs); err != nil {
				log.Fatal(err)
			}
			println(doc.String())
		case <-sigs:
			log.Println("Bye")
			return
		}
	}
}

func (d *Doc) InsertProjects(cl *gitlab.Client, projs []*gitlab.Project) error {
	for _, p := range projs {
		d.Insert(p)
		files, _ := gu.ListDir(cl, p.ID, *flgDir)
		for i := range files {
			log.Printf("Downloading %q %s\n", path.Join(p.WebURL, files[i].Path), files[i].Type)
			buf, err := gitlabutil.Download(cl, p.ID, "master", files[i].Path)
			if err != nil {
				return err
			}
			d.InsertFile(p, files[i].Path, buf)
		}
	}

	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return err
	}

	log.Println("Created index, now indexing")
	for _, g := range d.projects { // not thread safe
		for k, buf := range g.Files {
			if err := index.Index(k, string(buf)); err != nil {
				return err
			}
		}

	}
	d.SetIndex(index)
	count, err := index.DocCount()
	if err != nil {
		return err
	}
	log.Printf("Downloaded and indexed %d files, starting web server\n", count)
	return nil
}
