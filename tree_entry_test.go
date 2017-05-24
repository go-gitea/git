// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"io/ioutil"
	"os"
	"testing"
)

func setupGitRepo(url string) string {
	dir, err := ioutil.TempDir("", "gitea-bench")
	if err != nil {
		panic(err)
	}
	/* Manual method
	_, err = NewCommand("clone", url, dir).Run()
	if err != nil {
		log.Fatal(err)
	}
	*/
	err = Clone(url, dir, CloneRepoOptions{})
	if err != nil {
		panic(err)
	}
	return dir
}

//TODO use https://blog.golang.org/subtests when removing support for Go1.6
func benchmarkGetCommitsInfo(url string, b *testing.B) {
	b.StopTimer()
	
	// setup env
	repoPath := setupGitRepo(url)
	defer os.RemoveAll(repoPath)

	repo, err := OpenRepository(repoPath)
	if err != nil {
		panic(err)
	}

	commit, err := repo.GetBranchCommit("master")
	if err != nil {
		panic(err)
	}

	entries, err := commit.Tree.ListEntries()
	if err != nil {
		panic(err)
	}
	entries.Sort()

	b.StartTimer()
	// run the GetCommitsInfo function b.N times
	for n := 0; n < b.N; n++ {
		_, err = entries.GetCommitsInfo(commit, "")
		if err != nil {
			panic(err)
		}
	}
}


func BenchmarkGetCommitsInfoGitea(b *testing.B)  { benchmarkGetCommitsInfo("https://github.com/go-gitea/gitea.git", b) } //5k+ commits
func BenchmarkGetCommitsInfoMoby(b *testing.B)  { benchmarkGetCommitsInfo("https://github.com/moby/moby.git", b) } //32k+ commits
func BenchmarkGetCommitsInfoGo(b *testing.B)  { benchmarkGetCommitsInfo("https://github.com/golang/go.git", b) } //32k+ commits
func BenchmarkGetCommitsInfoLinux(b *testing.B)  { benchmarkGetCommitsInfo("https://github.com/torvalds/linux.git", b) } //677k+ commits
func BenchmarkGetCommitsInfoManyFile(b *testing.B)  { benchmarkGetCommitsInfo("https://github.com/ethantkoenig/manyfiles", b) }
