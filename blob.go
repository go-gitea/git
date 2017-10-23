// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"io"
)

// Blob represents a Git object.
type Blob struct {
	repo *Repository
	*TreeEntry
}

// Data gets content of blob all at once and wrap it as io.Reader.
// WARNING: the io.PipeReader should be read or close on the invoke place
func (b *Blob) Data() (*io.PipeReader, error) {
	r, w := io.Pipe()

	var err error
	go func() {
		defer w.Close()

		stderr := new(bytes.Buffer)
		err = b.DataPipeline(w, stderr)
		err = concatenateError(err, stderr.String())
	}()

	return r, err
}

// DataPipeline gets content of blob and write the result or error to stdout or stderr
func (b *Blob) DataPipeline(stdout, stderr io.Writer) error {
	return NewCommand("show", b.ID.String()).RunInDirPipeline(b.repo.Path, stdout, stderr)
}

// Bytes gets content of blob all at once.
// This can be very slow and memory consuming for huge content but for small object it's enough
func (b *Blob) Bytes() ([]byte, error) {
	return NewCommand("show", b.ID.String()).RunInDirBytes(b.repo.Path)
}
