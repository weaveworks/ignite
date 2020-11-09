#!/bin/bash

PREVIOUS_TAG=${1}
RELEASE_NOTES=${2}

usage() {
  echo "This script is useful for finding uncredited/undocumented PR's in a GREN generated release_notes/changelog file."
  echo "It checks the git merge log on your local machine for PR numbers that aren't present in the GREN changelog."
  echo
  echo "usage:"
  echo "  $0 <previous_tag> <release_notes_file>"
  echo
  echo "example:"
  echo "  check for PR's missing from the v0.8.0 release notes"
  echo
  echo "  $0 v0.7.1 ./docs/releases/v0.8.0.md"
}

if ! [ "${1}" ] || ! [ "${2}" ]; then
  usage
  exit 1
fi

recent_prs_from_git() {
  git log "${PREVIOUS_TAG}.." --merges --oneline | grep -o '#[1-9][0-9]*'
}

missing_prs() {
  for pr in $(recent_prs_from_git); do
    if ! grep "$pr" "${RELEASE_NOTES}" >/dev/null; then
      echo "$pr"
    fi
  done
}

format_grep_filter() {
  printf "|%s" "$@" | sed 's/^|//'
}

hub pr list -s closed \
  | grep -E "$(format_grep_filter $(missing_prs))"
