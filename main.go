package main

import (
	"fmt"
	"log"

	gu "github.com/miekg/xdoc/gitlabutil"
	"github.com/xanzy/go-gitlab"
)

const DocDir = "xdoc"

func main() {
	base := "https://gitlab.gnome.org"
	cl := gu.NewClient(base)
	groups, err := gu.ListGroups(cl)
	if err != nil {
		log.Fatal(err)
	}

	for i := range groups {
		fmt.Printf("%d %s %+v\n", groups[i].ID, groups[i].Name, groups[i].Projects)
		printProj(cl, groups[i].ID)
		sub, err := gu.ListSubgroups(cl, groups[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		for j := range sub {
			fmt.Printf("    %d %s\n", sub[j].ID, sub[j].Name)
			printProj(cl, sub[j].ID)
		}
	}
}

func printProj(cl *gitlab.Client, gid int) {
	proj, _ := gu.ListProjects(cl, gid)
	for i := range proj {
		fmt.Println("-- " + proj[i].Name + ": " + proj[i].WebURL)
	}
}
