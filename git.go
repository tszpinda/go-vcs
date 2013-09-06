package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

func (r *gitRepo) Log(startRev, endRev string) ([]Log, error) {
	arg := ""
	if startRev != "" {
		arg = startRev
		if endRev != "" {
			arg = startRev + ".." + endRev
		}
	}
	var cmd *exec.Cmd
	if arg == "" {
		cmd = exec.Command("git", "log", "--pretty=format:'%h|%an|%ad|%s'", "--date=short")
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:'%h|%an|%ad|%s'", "--date=short", arg)
	}
	cmd.Dir = r.dir
	if out, err := cmd.CombinedOutput(); err == nil {
		commits := strings.Split(strings.Replace(string(out), "'", "", -1), "\n")
		found := len(commits)
		//check if last element was \n so it its empty
		if found > 0 && commits[len(commits)-1] == "" {
			//remove last one
			commits = commits[:len(commits)-1]
		}

		logs := make([]Log, len(commits))
		for i, commit := range commits {
			commitArr := strings.Split(commit, "|")
			date, _ := time.Parse("2006-01-02", commitArr[2])
			log := Log{commitArr[0], commitArr[1], date, commitArr[3]}
			logs[i] = log
		}

		return logs, nil
	} else {
		if strings.Contains(string(out), "unknown revision") {
			return nil, fmt.Errorf("One or both revisions not found: '%v' - '%v'", startRev, endRev)
		}
		return nil, fmt.Errorf("git log failed:\n error details:\n%s\n%s", startRev, endRev, err, out)
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
