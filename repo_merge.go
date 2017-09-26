package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//Merge creates a new commit for the repo which is composed of a three-way merge of the branches.
func (repo *Repository) MergeCommit(sig *Signature, baseID, headID, message string) (string, error) {

	//
	//Create required temporary directory/file
	//
	workTreePath, err := ioutil.TempDir("", "git-merge-worktree-")
	if err != nil {
		return "", fmt.Errorf("Cannot create temporary directory git-merge-worktree-*: %v", err)
	}
	defer os.RemoveAll(workTreePath)

	f, err := ioutil.TempFile("", "git-merge-index-")
	if err != nil {
		return "", err
	}
	indexPath := f.Name()
	f.Close()
	defer os.Remove(indexPath)

	//
	//Construct the tree object from the input commits
	//
	envForReadTree := []string{
		"GIT_DIR=" + repo.Path,
		"GIT_WORK_TREE=" + workTreePath,
	}

	commonID, err := repo.GetMergeBase(baseID, headID)
	if err != nil {
		return "", err
	}
	if _, err = NewEnvCommand(envForReadTree, "read-tree", "--index-output", indexPath, "-im", commonID, baseID, headID).RunInDir(repo.Path); err != nil {
		return "", err
	}

	envForMergeIndex := []string{
		"GIT_DIR=" + repo.Path,
		"GIT_INDEX_FILE=" + indexPath,
		"GIT_WORK_TREE=" + workTreePath,
	}

	if _, err = NewEnvCommand(envForMergeIndex, "merge-index", "git-merge-one-file", "-a").RunInDir(repo.Path); err != nil {
		return "", err
	}

	envForWriteTree := envForMergeIndex
	_treeID, err := NewEnvCommand(envForWriteTree, "write-tree").RunInDir(repo.Path)
	if err != nil {
		return "", err
	}
	treeID := strings.TrimSpace(_treeID)

	//
	//Construct the merge commit object from the tree object we have created.
	//
	headCommit, err := repo.GetCommit(headID)
	if err != nil {
		return "", err
	}
	envForCommit := []string{
		"GIT_DIR=" + repo.Path,
		"GIT_AUTHOR_NAME=" + headCommit.Author.Name,
		"GIT_AUTHOR_EMAIL=" + headCommit.Author.Email,
		"GIT_AUTHOR_DATE=" + headCommit.Author.When.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		"GIT_COMMITTER_NAME=" + sig.Name,
		"GIT_COMMITTER_EMAIL=" + sig.Email,
		"GIT_COMMITTER_DATE=" + sig.When.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
	}

	_mergeCommit, err := NewEnvCommand(envForCommit, "commit-tree", treeID, "-p", baseID, "-p", headID, "-m", message).RunInDir(repo.Path)
	if err != nil {
		return "", err
	}
	mergeCommit := strings.TrimSpace(_mergeCommit)
	return mergeCommit, nil
}

//MergeBranch merges head into base with the commit message. It updates the ref of base as well.
func (repo *Repository) MergeBranch(sig *Signature, base, head, message string) error {
	oursID, err := repo.GetBranchCommitID(base)
	if err != nil {
		return err
	}
	theirsID, err := repo.GetBranchCommitID(head)
	if err != nil {
		return err
	}
	mergeCommitID, err := repo.MergeCommit(sig, oursID, theirsID, message)
	if err != nil {
		return err
	}

	env := []string{
		"GIT_DIR=" + repo.Path,
	}
	ref := fmt.Sprintf("refs/heads/%s", base)
	if _, err := NewEnvCommand(env, "update-ref", ref, mergeCommitID).RunInDir(repo.Path); err != nil {
		return err
	}

	if _, err := NewEnvCommand(env, "pack-refs", "--all").RunInDir(repo.Path); err != nil {
		return err
	}

	return nil
}
