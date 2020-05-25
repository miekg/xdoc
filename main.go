package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/miekg/xdoc/gitlabutil"
	gu "github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
)

const DocDir = "docs"

var (
	flgGroup  = flag.String("group", "", "select only this group")
	flgBase   = flag.String("base", "", "base URL to add")
	flgGitlab = flag.String("gitlab", "https://gitlab.com", "GitLab site")
	flgDir    = flag.String("dir", DocDir, "directory to use for documentation")
	flgInt    = flag.Duration("int", 10*time.Minute, "duration to sleep before restarting the download loop")
	flgWorker = flag.Uint("worker", 10, "number of concurrent worker to download from gitlab")
)

func main() {
	flag.Parse()
	if *flgGroup == "" {
		log.Fatal("-group need a value")
	}

	cl := gu.NewClient(*flgGitlab)
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

type downloadItem struct {
	project  *gitlab.Project
	ref      string
	filename string
}

type downloadResult struct {
	project  *gitlab.Project
	filename string
	buf      []byte
	err      error
}

func downloadWorker(cl *gitlab.Client, work <-chan downloadItem, result chan<- downloadResult) {
	for w := range work {
		buf, err := gitlabutil.Download(cl, w.project.ID, w.ref, w.filename)
		result <- downloadResult{w.project, w.filename, buf, err}
	}
	log.Printf("Goroutine download worker shutting down")
}

func (d *Doc) InsertProjects(cl *gitlab.Client, projs []*gitlab.Project) error {

	work := make(chan downloadItem, *flgWorker)
	result := make(chan downloadResult, *flgWorker)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < int(*flgWorker); i++ {
		go downloadWorker(cl, work, result)
	}

	go func() {
		for {
			select {
			case r := <-result:
				if r.err != nil {
					log.Printf("Error when downloading %q: %s", r.filename, r.err)
					continue
				}
				d.InsertFile(r.project, r.filename, r.buf)
				wg.Done()

			case <-stop:
				return
			}
		}
	}()

	for _, p := range projs {
		opts, _ := ParseOptions(cl, p)
		d.Insert(p, opts)
		files, _ := gu.ListDir(cl, p.ID, *flgDir)
		wg.Add(len(files))
		for i := range files {
			log.Printf("Downloading %q %s\n", path.Join(p.WebURL, files[i].Path), files[i].Type)
			work <- downloadItem{p, "master", files[i].Path}
		}
	}
	close(work)
	wg.Wait()
	close(stop)

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
