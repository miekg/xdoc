package main

import (
	"log"

	"github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type Options struct {
	Ref    string `yaml:"ref"`
	Lang   string `yaml:"lang"`
	Flavor `yaml:"flavor"`
}

// ParseOptions parses the .<doc>.yaml options file in the root of the repository in the master branch
func ParseOptions(cl *gitlab.Client, p *gitlab.Project) (Options, error) {
	buf, err := gitlabutil.Download(cl, p.ID, "master", *flgDir+".yaml")
	if err != nil {
		return Options{}, err
	}
	log.Printf("Project %s has options YAML file", p.WebURL)
	opt := Options{}
	if err := yaml.Unmarshal(buf, &opt); err != nil {
		return opt, err
	}
	if opt.Ref == "" {
		opt.Ref = "master"
	}
	return opt, nil
}
