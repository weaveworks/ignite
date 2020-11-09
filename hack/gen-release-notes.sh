#!/bin/bash

PREVIOUS_TAG=${1}
OPTIONAL_REF=${2}

usage() {
  echo "This script generates a starter release note document using `git` and `hub`."
  echo "It finds all of the PR's with matching commits in the git log since the indicated ref."
  echo "The PR's are link formatted with their title and author and are sorted by label."
  echo
  echo "usage:"
  echo "  $0 <previous_tag>"
  echo "  $0 <ref1> <ref2>"
  echo
  echo "examples:"
  echo
  echo "  $0 v0.7.1"
  echo "  $0 v0.6.0 v0.7.0"
}

if ! [ "${1}" ]; then
  usage
  exit 1
fi

gh_prs() {
  hub pr list -s closed \
    --format "_commits_%sm_%sB_%sH_commits_ %l - %t ([%i](%U), [%au](https://github.com/%au))%n"
}

git_shas_since() {
  git log --format="_%H_" "$1..$2"
}

grep_filter() {
  printf "|%s" "$@" | sed 's/^|//'
}

gh_prs \
  | grep -E "$(grep_filter $(git_shas_since ${PREVIOUS_TAG} ${OPTIONAL_REF}))" \
  | sed 's/^_commits_[^ ]*_commits_//' \
  | sort -r  `# reverse sort puts all of the un-labelled PR's at the bottom`

