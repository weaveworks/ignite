#!/usr/bin/env bash

PREVIOUS_TAG=${1}
RELEASE_NOTES=${2}

usage() {
  echo "This script is useful for finding extra PR's in a GREN generated release_notes/changelog file."
  echo "It checks the GREN changelog for PR numbers that aren't present the git merge log on your local machine."
  echo
  echo "usage:"
  echo "  $0 <previous_tag> <release_notes_file>"
  echo
  echo "example:"
  echo "  check for extra PR's in the v0.8.0 release notes"
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

release_note_prs() {
  grep -o -E "pull/[0-9][0-9]*" "${RELEASE_NOTES}" | sed 's,pull/,#,'
}

format_grep_filter() {
  printf "|%s" "$@" | sed 's/^|//'
}

missing_prs() {
  release_note_prs | grep -v -E "$(format_grep_filter $(recent_prs_from_git))"
}

hub pr list -s closed \
  | grep -E "$(format_grep_filter $(missing_prs))"
