package main

import (
	"net/url"
	"sync"

	"github.com/xanzy/go-gitlab"
)

type Doc struct {
	Loc      string            // Where is the markdown stored
	Projects map[string]GitLab // Basename(URL) -> Project + potential metadata
	rw       sync.RWMutex      // protects Projects
}

type GitLab struct {
	*gitlab.Project
	Commit string // not used, here for future expansion
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
