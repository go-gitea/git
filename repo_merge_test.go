package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRepoMerge(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "git-testing-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	for _, x := range []struct {
		//Input
		repoName     string
		oursBranch   string
		theirsBranch string
		author       *Signature
		committer    *Signature
		message      string

		//Output
		mergeCommitID string
	}{
		{
			repoName:      "repo-ff",
			oursBranch:    "master",
			theirsBranch:  "bye-world",
			author:        &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 20, 0, 0, 0, 0, time.UTC)},
			committer:     &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 25, 0, 0, 0, 0, time.UTC)},
			message:       "Merge branch 'bye-world' of repo-ff",
			mergeCommitID: "b4410a0fa606b697c892a3de8bbced3c005422a0",
		},
	} {
		from := filepath.Join("testdata", x.repoName)
		to := filepath.Join(tmpdir, x.repoName)
		if err := Clone(from, to, CloneRepoOptions{
			Timeout: 5 * time.Second,
			Mirror:  true,
			Bare:    true,
			Quiet:   true,
			Branch:  "",
		}); err != nil {
			t.Fatal(err)
		}

		repo, err := OpenRepository(to)
		if err != nil {
			t.Fatal(err)
		}

		if err := repo.MergeBranch(x.committer, x.oursBranch, x.theirsBranch, x.message); err != nil {
			t.Fatal(err)
		}

		oursID, err := repo.GetBranchCommitID(x.oursBranch)
		if err != nil {
			t.Fatal(err)
		}

		if oursID != x.mergeCommitID {
			t.Fatalf("Not matched. Actual:%s Expected:%s", oursID, x.mergeCommitID)
		}
	}
}
