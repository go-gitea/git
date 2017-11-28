// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readPackedRefs(repoPath string) ([]string, error) {
	path := repoPath + "/packed-refs"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	refData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("v\\d+\\.\\d+\\.\\d+$")
	names := []string{}

	for _, ref := range bytes.Split(refData, []byte("\n")) {
		if tag := re.Find(ref); tag != nil {
			names = append(names, string(tag))
		}
	}

	return names, nil
}

func readRefDir(repoPath, prefix, relPath string) ([]string, error) {
	dirPath := filepath.Join(repoPath, prefix, relPath)
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(fis))
	for _, fi := range fis {
		if strings.Contains(fi.Name(), ".DS_Store") {
			continue
		}

		relFileName := filepath.Join(relPath, fi.Name())
		if fi.IsDir() {
			subnames, err := readRefDir(repoPath, prefix, relFileName)
			if err != nil {
				return nil, err
			}
			names = append(names, subnames...)
			continue
		}

		names = append(names, relFileName)
	}

	return names, nil
}
