package vcs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	t.Parallel()

	url := "https://bitbucket.org/pawmar/golang_playground"
	//tmpdir := "/Users/tszpinda/tmp/git-test"
	var tmpdir string
	tmpdir, err := ioutil.TempDir("", "go-vcs-TestGit")
	if err != nil {
		t.Fatalf("TempDir: %s", err)
	}
	tmpdir = "/Users/tszpinda/tmp/git-test"
	//defer os.RemoveAll(tmpdir)

	r, err := CloneOrOpen(Git, url, tmpdir)
	if err != nil {
		t.Fatalf("Error CloneOrOpen: %s", err)
	}
	fmt.Println("Checkout master")

	masterDir, err := r.CheckOut("master")
	if err != nil {
		t.Fatalf("CheckOut master: %s", err)
	}
	fmt.Println("master dir", masterDir)

	fmt.Println("Retrive logs")
	r.Pull()
	//in Git range is exclusive from bottom, inclusive from top
	logs, err := r.Log("2dbee0b", "33422dc")
	if err != nil {
		t.Fatalf("Log error: %s", err)
	}

	if len(logs) != 3 {
		fmt.Printf("Expected 3 but found: %v", len(logs))
		t.Fail()
	}
	if logs[0].Hash != "33422dc" {
		fmt.Println("Invalid hash")
		t.Fail()
	}
	if logs[0].User != "Pawel Markowski" {
		fmt.Println("Invalid hash")
		t.Fail()
	}
	if logs[0].Message != "no staged" {
		fmt.Println("Invalid hash")
		t.Fail()
	}

	//invalid rev
	logs, err = r.Log("2dbee0b", "33422dk")
	if logs != nil {
		t.Fatalf("Logs should be nil")
	}
	if err == nil {
		t.Fatalf("Err should be not be empty")
	}
	expErr := fmt.Errorf("One or both revisions not found: '%v' - '%v'", "2dbee0b", "33422dk")
	if err.Error() != expErr.Error() {
		t.Fatalf("Err should be :'%v' but was: '%v'", expErr, err)
	}
}

func TestGit(t *testing.T) {
	t.Parallel()

	var tmpdir string
	tmpdir, err := ioutil.TempDir("", "go-vcs-TestGit")
	if err != nil {
		t.Fatalf("TempDir: %s", err)
	}
	defer os.RemoveAll(tmpdir)

	url := "https://bitbucket.org/sqs/go-vcs-gittest.git"
	r, err := Clone(Git, url, tmpdir)
	if err != nil {
		t.Fatalf("Clone: %s", err)
	}

	// check out master
	masterDir, err := r.CheckOut("master")
	if err != nil {
		t.Fatalf("CheckOut master: %s", err)
	}
	assertFileContains(t, masterDir, "foo", "Hello, foo\n")
	assertNotFileExists(t, masterDir, "bar")

	// check out a branch
	barbranchDir, err := r.CheckOut("barbranch")
	if err != nil {
		t.Fatalf("CheckOut barbranch: %s", err)
	}
	assertFileContains(t, barbranchDir, "bar", "Hello, bar\n")

	// check out a commit id
	barcommit := "f411e1ea59ed2b833291efa196e8dab80dbf7cb8"
	barcommitDir, err := r.CheckOut(barcommit)
	if err != nil {
		t.Fatalf("CheckOut barcommit %s: %s", barcommit, err)
	}
	assertFileContains(t, barcommitDir, "bar", "Hello, bar\n")

	if _, err := Clone(Git, url, tmpdir); !os.IsExist(err) {
		t.Fatalf("Clone to existing dir: want os.IsExist(err), got %T %v", err, err)
	}

	// Open
	if r, err = Open(Git, tmpdir); err != nil {
		t.Fatalf("Open: %s", err)
	}
	if masterDir, err = r.CheckOut("master"); err != nil {
		t.Fatalf("CheckOut master: %s", err)
	}
	assertFileContains(t, masterDir, "foo", "Hello, foo\n")
}
