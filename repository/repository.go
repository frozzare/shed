package repository

import (
	"errors"
	"os"
	"strings"

	"github.com/frozzare/shed/config"
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

// NewRepository creates a new repository.
func NewRepository(config config.Git) (Repository, error) {
	var err error
	var path string

	if len(config.Path) > 0 {
		path = config.Path
	} else {
		path, err = os.Getwd()
	}

	if err != nil {
		return Repository{}, err
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		if len(config.Branch) > 0 {
			return NewRepositoryFromBranch(config.Branch), nil
		}

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

	return NewRepositoryFromBranch(branch), nil
}

// NewRepositoryFromBranch creates a new repository from branch name.
func NewRepositoryFromBranch(branch string) Repository {
	return Repository{
		Branch: branch,
		Slug:   slug.Make(branch),
	}
}
