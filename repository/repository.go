package repository

import (
	"errors"
	"os"
	"strings"

	"github.com/gosimple/slug"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	ErrNoBranch = errors.New("Git head is not a branch")
)

// Repository represents the git repository.
type Repository struct {
	Branch string
	Slug   string
}

// NewRepository creates a new repostiory.
func NewRepository(args ...string) (Repository, error) {
	var err error
	var path string

	if len(args) > 0 && args[0] != "" {
		path = args[0]
	} else {
		path, err = os.Getwd()
	}

	if err != nil {
		return Repository{}, err
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return Repository{}, err
	}

	var branch string

	h, err := r.Head()

	if err == plumbing.ErrReferenceNotFound {
		h, err = r.Reference(plumbing.HEAD, false)

		if err != nil {
			return Repository{}, err
		}

		n := h.Target().String()
		p := strings.Split(n, "/")
		branch = p[len(p)-1]
	} else if err != nil {
		return Repository{}, err
	} else if !h.IsBranch() {
		return Repository{}, ErrNoBranch
	}

	if len(branch) == 0 {
		n := h.Name().String()
		p := strings.Split(n, "/")
		branch = p[len(p)-1]
	}

	return Repository{
		Branch: branch,
		Slug:   slug.Make(branch),
	}, nil
}
