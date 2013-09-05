package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type git struct {
	cmd string
}

var Git VCS = git{"git"}

type gitRepo struct {
	dir string
	git *git
}

func (git git) Clone(url, dir string) (Repository, error) {
	r := &gitRepo{dir, &git}

	cmd := exec.Command("git", "clone", "--", url, dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		if strings.Contains(string(out), fmt.Sprintf("fatal: destination path '%s' already exists", dir)) {
			return nil, os.ErrExist
		}
		return nil, fmt.Errorf("git clone failed: %s\n%s", err, out)
	}

	return r, nil
}

func (git git) Open(dir string) (Repository, error) {
	// TODO(sqs): check for .git or bare repo
	if _, err := os.Stat(dir); err == nil {
		return &gitRepo{dir, &git}, nil
	} else {
		return nil, err
	}
}

func (r *gitRepo) Dir() (dir string) {
	return r.dir
}

func (r *gitRepo) VCS() VCS {
	return r.git
}

func (r *gitRepo) Download() error {
	panic("not implemented")
}

func (r *gitRepo) CheckOut(rev string) (dir string, err error) {
	cmd := exec.Command("git", "checkout", rev)
	cmd.Dir = r.dir
	if out, err := cmd.CombinedOutput(); err == nil {
		return r.dir, nil
	} else {
		return "", fmt.Errorf("git checkout %q failed: %s\n%s", rev, err, out)
	}
}

func (r *gitRepo) Log(startRev, endRev string) ([]string, error) {
	arg := startRev + ".." + endRev
	cmd := exec.Command("git", "log", "--pretty=oneline", "--abbrev-commit", arg)
	cmd.Dir = r.dir
	if out, err := cmd.CombinedOutput(); err == nil {
		log := string(out)
		logs := strings.Split(log, "\n")
		found := len(logs)
		//check if last element was \n so it its empty
		if found > 0 && logs[len(logs)-1] == "" {
			//remove last one
			logs = logs[:len(logs)-1]
		}

		return logs, nil
	} else {
		return nil, fmt.Errorf("git log --pretty=oneline --abbrev-commit %v..%v failed:\n error details:\n%s\n%s", startRev, endRev, err, out)
	}
}

func (r *gitRepo) HardReset() error {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = r.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch failed: %s\n%s", err, out)
	}

	cmd = exec.Command("git", "reset", "--hard", "origin/master")
	cmd.Dir = r.dir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git reset --hard origin/master: %s\n%s", err, out)
	}
	return nil
}

func (r *gitRepo) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = r.dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %s\n%s", err, out)
	}
	return nil
}
