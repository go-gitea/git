// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	version "github.com/mcuadros/go-version"
)

// Storage is an interface which describes how to fetch git objects from the git data storage placement.
type Storage interface {
	GetTags(string) ([]string, error)
}

// shellStorage is an implementation of Storage which use git shell to retrieve git objects from file system.
type shellStorage struct {
}

// GetTags return the tag names according repoPath
func (shellStorage) GetTags(repoPath string) ([]string, error) {
	cmd := NewCommand("tag", "-l")
	if version.Compare(gitVersion, "2.0.0", ">=") {
		cmd.AddArguments("--sort=-v:refname")
	}

	stdout, err := cmd.RunInDir(repoPath)
	if err != nil {
		return nil, err
	}

	tags := strings.Split(stdout, "\n")
	tags = tags[:len(tags)-1]

	if version.Compare(gitVersion, "2.0.0", "<") {
		version.Sort(tags)

		// Reverse order
		for i := 0; i < len(tags)/2; i++ {
			j := len(tags) - i - 1
			tags[i], tags[j] = tags[j], tags[i]
		}
	}

	return tags, nil
}

type localStorage struct {
}

// GetTags return the tag names according repoPath
func (localStorage) GetTags(repoPath string) ([]string, error) {
	packed, err := readPackedRefs(repoPath)
	if err != nil {
		return nil, err
	}

	if !com.IsExist(filepath.Join(repoPath, "refs/tags")) {
		return packed, nil
	}

	// Attempt loose files first as the /refs/tags folder should always
	// exist whether it has files or not.
	loose, err := readRefDir(repoPath, "refs/tags", "")
	if err != nil {
		return nil, err
	}

	// If both loose refs and packed refs exist then it's highly
	// likely that the loose refs are more recent than packed (created
	// on top of packed older refs). Therefore we can append each
	// together taking the packed refs first.
	return append(packed, loose...), nil
}

var (
	// ShellStorage provides methods to retrieve git objects via git shell
	ShellStorage = new(shellStorage)
	// LocalStorage provides methods to retrieve git objects via pure go
	LocalStorage = new(localStorage)
	// DefaultStorage is the default implementation of git data storage
	DefaultStorage Storage = ShellStorage
)
