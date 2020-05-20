package gitlabutil

import "github.com/xanzy/go-gitlab"

func ListDir(cl *gitlab.Client, pid int, dir string) ([]*gitlab.TreeNode, error) {
	listopts := gitlab.ListOptions{PerPage: 50, Page: 1}
	opts := &gitlab.ListTreeOptions{ListOptions: listopts, Path: gitlab.String(dir)}

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

// Download will download a file from gitlab.
func Download(cl *gitlab.Client, pid int, name string) ([]byte, error) {
	opts := &gitlab.GetRawFileOptions{}

	data, _, err := cl.RepositoryFiles.GetRawFile(pid, name, opts)
	return data, err
}
