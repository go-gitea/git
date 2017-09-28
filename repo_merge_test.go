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
		errorString   string
	}{
		{
			repoName:      "repo-ff",
			oursBranch:    "master",
			theirsBranch:  "bye-world",
			author:        &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 20, 0, 0, 0, 0, time.UTC)},
			committer:     &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 25, 0, 0, 0, 0, time.UTC)},
			message:       "Merge branch 'bye-world' of repo-ff",
			errorString:   "",
			mergeCommitID: "0896984c31681a77c7660bacd56de91c15389e88",
		},
		{
			repoName:      "repo-conflict",
			oursBranch:    "bye-world",
			theirsBranch:  "greetings",
			author:        &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 20, 0, 0, 0, 0, time.UTC)},
			committer:     &Signature{Email: "Alice@example.com", Name: "Alice", When: time.Date(2017, time.September, 25, 0, 0, 0, 0, time.UTC)},
			message:       "Merge branch 'greetings' of repo-conflict",
			errorString:   "exit status 128 - ERROR: content conflict in README.md\nfatal: merge program failed\n",
			mergeCommitID: "",
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
			if err.Error() != x.errorString {
				t.Fatal(err)
			}
		} else {
			oursID, err := repo.GetBranchCommitID(x.oursBranch)
			if err != nil {
				t.Fatal(err)
			}

			if oursID != x.mergeCommitID {
				t.Fatalf("Not matched. Actual:%s Expected:%s", oursID, x.mergeCommitID)
			}
		}
	}
}
