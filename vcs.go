package vcs

import (
	"os"
	"time"
)

type VCS interface {
	// Clones the repository at the given URL into dir. If dir already exists, the error os.ErrExist
	// is returned.
	Clone(url, dir string) (Repository, error)

	Open(dir string) (Repository, error)
}

// Map of VCS name to VCS object.
var VCSByName = map[string]VCS{
	"git": Git,
	"hg":  Hg,
}

type Log struct {
	Hash    string
	User    string
	Date    time.Time
	Message string
}

type Repository interface {
	Dir() string // The repository's root directory.
	VCS() VCS

	// Downloads updates to the repository from the default remote.
	Download() error

	// CheckOut returns the path of a directory containing a working tree at revision rev. CheckOut
	// assumes that rev is local or has already been fetched; it does not update the repository.
	CheckOut(rev string) (dir string, err error)

	// Retruns list of logs from the commits between two revisions
	Log(startRev, endRev string) ([]Log, error)

	//Hard rest
	HardReset() error

	//Pull
	Pull() error
}

// Clones the VCS repository from a remote URL to dir.
func Clone(vcs VCS, url, dir string) (Repository, error) {
	return vcs.Clone(url, dir)
}

// Opens the VCS repository at dir.
func Open(vcs VCS, dir string) (Repository, error) {
	return vcs.Open(dir)
}

// If no repository exists at dir, CloneOrOpen clones the VCS repository to dir. Otherwise, it opens
// the repository at dir (without checking that the repository there is, indeed, cloned from the
// specified URL).
func CloneOrOpen(vcs VCS, url, dir string) (Repository, error) {
	if repo, err := Clone(vcs, url, dir); os.IsExist(err) {
		return Open(vcs, dir)
	} else {
		return repo, err
	}
}
