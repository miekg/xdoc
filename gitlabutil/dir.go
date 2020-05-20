package gitlabutil

import "github.com/xanzy/go-gitlab"

func ListDir(cl *gitlab.Client, pid int) ([]*gitlab.TreeNode, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListTreeOptions{ListOptions: listopts, Path: gitlab.String("xdoc")}

	trees := []*gitlab.TreeNode{}
	for {
		tree, resp, err := cl.Repositories.ListTree(pid, opts)
		if err != nil {
			return trees, err
		}
		trees = append(trees, tree...)
		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	return trees, nil
}
