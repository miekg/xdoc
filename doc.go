package main

import (
	"net/url"
	"path"
	"sync"

	"github.com/xanzy/go-gitlab"
)

// Flavor is the markdown flavor used for parsing the markdown files.
type Flavor int

const (
	None Flavor = iota
	Mmark
)

type Doc struct {
	Loc      string            // Where is the markdown stored
	Projects map[string]GitLab // Basename(URL) -> Project + potential metadata
	rw       sync.RWMutex      // protects Projects
}

type GitLab struct {
	*gitlab.Project
	Commit string // not used yet
	Lang   string // not used yet
	Flavor        // not used yet
}

// Init initializes d.
func (d *Doc) Init() {
	// create tmp dir for where we can download the markdown.
}

// Insert inserts a new project into d.
func (d *Doc) Insert(p *gitlab.Project) {
	d.rw.Lock()
	defer d.rw.Unlock()
	url, _ := url.Parse(p.WebURL)
	d.Projects[url.Path] = GitLab{Project: p, Commit: ""}
}

// Fetch will return the project belonging to path. Will return nil if not found.
func (d *Doc) Fetch(path string) GitLab {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.Projects[path]
}

// FullPath returns the on-disk path for this gitlab project and path.
func (d *Doc) FullPath(g GitLab, pathname string) string {
	a := path.Join(d.Loc, ProjectToPath(g.Project))
	b := path.Join(a, pathname)
	return b
}

// ProjectToPath converts a gitlab project to a path that can be used in Fetch.
func ProjectToPath(p *gitlab.Project) string {
	url, _ := url.Parse(p.WebURL)
	return url.Path
}
