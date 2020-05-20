package gitlabutil

import (
	"github.com/xanzy/go-gitlab"
)

func NewClient(base string) *gitlab.Client {
	cl, _ := gitlab.NewClient("", gitlab.WithBaseURL(base))
	return cl
}
