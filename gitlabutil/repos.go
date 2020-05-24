package gitlabutil

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func ListGroups(cl *gitlab.Client) ([]*gitlab.Group, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListGroupsOptions{AllAvailable: gitlab.Bool(true), ListOptions: listopts}

	groups := []*gitlab.Group{}
	for {
		grp, resp, err := cl.Groups.ListGroups(opts)
		if err != nil {
			return groups, err
		}
		groups = append(groups, grp...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	return groups, nil
}

func GetGroup(cl *gitlab.Client, group string) ([]*gitlab.Group, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListGroupsOptions{Search: gitlab.String(group), ListOptions: listopts}

	groups := []*gitlab.Group{}
	for {
		grp, resp, err := cl.Groups.ListGroups(opts)
		if err != nil {
			return groups, err
		}
		groups = append(groups, grp...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("no groups found")
	}
	if len(groups) != 1 {
		return nil, fmt.Errorf("multiple groups found")
	}

	grp, _, err := cl.Groups.GetGroup(groups[0].ID)
	return []*gitlab.Group{grp}, err
}

func ListSubgroups(cl *gitlab.Client, gid int) ([]*gitlab.Group, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListSubgroupsOptions{ListOptions: listopts}

	groups := []*gitlab.Group{}
	for {
		grp, resp, err := cl.Groups.ListSubgroups(gid, opts)
		if err != nil {
			return groups, err
		}
		groups = append(groups, grp...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	return groups, nil
}

func ListProjects(cl *gitlab.Client, gid int) ([]*gitlab.Project, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListGroupProjectsOptions{ListOptions: listopts, Archived: gitlab.Bool(false)}

	projs := []*gitlab.Project{}
	for {
		proj, resp, err := cl.Groups.ListGroupProjects(gid, opts)
		if err != nil {
			return projs, err
		}
		projs = append(projs, proj...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	return projs, nil
}
