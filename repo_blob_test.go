// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetBlob(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	bareRepo1, err := OpenRepository(bareRepo1Path)
	assert.NoError(t, err)

	testCases := []struct {
		OID  string
		Data []byte
	}{
		{"e2129701f1a4d54dc44f03c93bca0a2aec7c5449", []byte("file1\n")},
		{"6c493ff740f9380390d5c9ddef4af18697ac9375", []byte("file2\n")},
		{"b1fc9917b618c924cf4aa421dae74e8bf9b556d3", []byte("Hi\n")},
	}

	for _, testCase := range testCases {
		blob, err := bareRepo1.GetBlob(testCase.OID)
		assert.NoError(t, err)

		dataReader, err := blob.Data()
		assert.NoError(t, err)

		data, err := ioutil.ReadAll(dataReader)
		assert.NoError(t, err)
		assert.Equal(t, testCase.Data, data)
	}
}
