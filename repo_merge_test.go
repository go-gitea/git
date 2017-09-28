package git

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestRepoMerge(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "git-testing-")
	if err != nil {
		t.Fatal(err)
	}
	//defer os.RemoveAll(tmpdir)

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
			repoName:      "ff.git",
			oursBranch:    "master",
			theirsBranch:  "test-branch",
			author:        &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			committer:     &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			message:       "Merge branch 'test-branch'",
			errorString:   "",
			mergeCommitID: "57b5d2ba9bceb08e8111c676418ddcc872c50cf2",
		},
		{
			repoName:      "conflict.git",
			oursBranch:    "test-branch1",
			theirsBranch:  "test-branch2",
			author:        &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			committer:     &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			message:       "Merge branch 'test-branch2'",
			errorString:   "exit status 128 - ERROR: content conflict in README.md\nfatal: merge program failed\n",
			mergeCommitID: "",
		},
		{
			repoName:      "adding.git",
			oursBranch:    "master",
			theirsBranch:  "test-branch",
			author:        &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			committer:     &Signature{Email: "alice@example.com", Name: "Alice", When: time.Date(2005, time.April, 7, 22, 13, 13, 0, time.UTC)},
			message:       "Merge branch 'test-branch'",
			errorString:   "",
			mergeCommitID: "7f19d72c899e7d8fd2d46724bc95dfd4904704b4",
		},
	} {
		from := filepath.Join("testdata/generated", x.repoName)
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
				t.Fatalf("'%s': %v", x.repoName, err)
			}
		} else {
			oursID, err := repo.GetBranchCommitID(x.oursBranch)
			if err != nil {
				t.Fatal(err)
			}

			if oursID != x.mergeCommitID {
				t.Fatalf("'%s': Not matched. Actual:%s Expected:%s", x.repoName, oursID, x.mergeCommitID)
			}
		}
	}
}
