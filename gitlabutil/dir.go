package gitlabutil

import (
	"github.com/xanzy/go-gitlab"
)

func ListDir(cl *gitlab.Client, pid int, dir string) ([]*gitlab.TreeNode, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListTreeOptions{ListOptions: listopts, Path: gitlab.String(dir)}

	trees := []*gitlab.TreeNode{}
	for {
		tree, resp, err := cl.Repositories.ListTree(pid, opts)
		if err != nil {
			return trees, err
		}
		for i := range tree {
			if tree[i].Type == "blob" {
				trees = append(trees, tree[i])
			}
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}
		opts.Page = resp.NextPage
	}
	return trees, nil
}

// Download will download a file from gitlab.
func Download(cl *gitlab.Client, pid int, ref, pathname string) ([]byte, error) {
	opts := &gitlab.GetRawFileOptions{Ref: gitlab.String(ref)}

	data, _, err := cl.RepositoryFiles.GetRawFile(pid, pathname, opts)
	return data, err
}
