#!/bin/bash

set -e
set +x

export GIT_AUTHOR_NAME="Alice"
export GIT_AUTHOR_EMAIL="alice@example.com"
export GIT_AUTHOR_DATE="Thu, 07 Apr 2005 22:13:13 +0000"

export GIT_COMMITTER_NAME="Alice"
export GIT_COMMITTER_EMAIL="alice@example.com"
export GIT_COMMITTER_DATE="Thu, 07 Apr 2005 22:13:13 +0000"

export TZ="UTC"

function create_repo {
	local readonly name=$1
	local readonly filename=$2
	local readonly content=$3
	local readonly message=$4

	mkdir $name
	cd $name
		git init --q
		echo "$content" > $filename 
		git add .
		git commit -q -m "$message"
	cd ..
}

function create_branch_and_modify {
	local readonly name=$1
	local readonly oldbranch=$2
	local readonly newbranch=$3
	local readonly action=$4

	cd  $name
		git checkout -q -b $newbranch $oldbranch
		eval "$action"
		git add .
		git commit -q -m "$action"
	cd ..
}

function test_merge {
	local readonly name=$1
	local readonly ours=$2
	local readonly theirs=$3
	local readonly options=$4

	local readonly cloned=$name.cloned

	git clone -q $name.git $cloned
	cd $cloned
		git checkout -q $ours
		git merge -q $options $theirs
		echo "$name -> `git rev-list --max-count=1 HEAD`"
	cd ..
	rm -rf $cloned
}

function sync_and_cleanup {
	local readonly name=$1

	git clone -q --mirror $name $name.git
	rm -rf ./$name
}

function repo-ff {
	local readonly repo_name="ff"
	local readonly action="sed -i 's/Hello/Bye/' README.md"
	create_repo $repo_name "README.md" "Hello, World" "Initial Commit"
	create_branch_and_modify $repo_name "master" "test-branch" "$action"
	sync_and_cleanup $repo_name
	test_merge $repo_name "master" "test-branch" "--no-edit --no-ff"
}

function repo-conflict {
	local readonly repo_name="conflict"
	local readonly action1="sed -i 's/Hello/Bye/' README.md"
	local readonly action2="sed -i 's/Hello/Greetings/' README.md"
	create_repo $repo_name "README.md" "Hello, World" "Initial Commit"
	create_branch_and_modify $repo_name "master" "test-branch1" "$action1"
	create_branch_and_modify $repo_name "master" "test-branch2" "$action2"
	sync_and_cleanup $repo_name
	set +e
	test_merge $repo_name "test-branch1" "test-branch2" "--no-edit --no-ff"
	set -e
}

function repo-adding {
	local readonly repo_name="adding"
	local readonly action="sed -i 's/Hello/Bye/' README.md && echo 0123456789 > a.txt"
	create_repo $repo_name "README.md" "Hello, World" "Initial Commit"
	create_branch_and_modify $repo_name "master" "test-branch" "$action"
	sync_and_cleanup $repo_name
	test_merge $repo_name "master" "test-branch" "--no-edit --no-ff"
}

function main() {
	rm -rf generated
	mkdir generated
	cd generated
		repo-ff
		repo-conflict
		repo-adding
	cd ..
}

main
