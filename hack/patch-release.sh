#!/usr/bin/env bash

SCRIPT_DIR=$( dirname "${BASH_SOURCE[0]}" )
cd ${SCRIPT_DIR}/..

source hack/release-common.sh

FORCE=${FORCE:-0}
GENERATED_GIT_VERSION=$(hack/ldflags.sh --version-only)
CURRENT_BRANCH=$(git branch | grep \* | cut -d ' ' -f2)

if [[ ${GENERATED_GIT_VERSION} =~ "-dirty" && ${FORCE} == 0 ]]; then
    echo "Won't try to do a release when the git state is dirty"
    exit 1
fi

if [[ (! ${CURRENT_BRANCH} =~ "release-") && ${FORCE} == 0 ]]; then
    echo "Won't try to do a release from a non-release branch"
    exit 1
fi

MAJOR=0
MINOR=${MINOR:-0}
PATCH=${PATCH:-0}
VERSION="v${MAJOR}.${MINOR}.${PATCH}"
RELEASE_BRANCH="release-${MAJOR}.${MINOR}"
EXTRA=${EXTRA:-""}
FULL_VERSION=${VERSION}${EXTRA}

if [[ ${MINOR} == "0" ]]; then
    echo "MINOR is mandatory"
    exit 1
fi

if [[ ${PATCH} == "0" ]]; then
    echo "PATCH is mandatory"
    exit 1
fi

# Calculate from where to build the changelog. For a prerelease, the changelog is incremental, but the
# changelog for a stable release includes notes from the last stable release on this release branch
if [[ ${EXTRA} == "" ]]; then
    PREVIOUS_TAG=$(git tag --sort taggerdate | sort -r | grep -E "v[0-9]+\.${MINOR}\.[0-9]+$" | head -1)
else
    PREVIOUS_TAG=$(git tag --sort taggerdate | sort -r | grep -E "v[0-9]+\.${MINOR}" | head -1)
fi

echo "Releasing version ${FULL_VERSION}"

tag_release() {
    read -p "Are you sure you want to tag the release ${FULL_VERSION}? [y/N] " confirm
    if [[ ! ${confirm} =~ ^[Yy]$ ]]; then
        exit 1
    fi

    git tag -f ${FULL_VERSION}
}

push_artifacts() {
    read -p "Are you sure you want to push the release ${FULL_VERSION} artifacts? [y/N] " confirm
    if [[ ! ${confirm} =~ ^[Yy]$ ]]; then
        cat <<- EOF
		Done! Next, do this:

		git push upstream --tags -f
		git push upstream ${RELEASE_BRANCH}
		EOF
        exit 1
    fi
    git push upstream --tags -f
    git push upstream ${RELEASE_BRANCH}
}

if [[ $1 == "tidy" ]]; then
    make_tidy_autogen
elif [[ $1 == "changelog" ]]; then
    write_changelog
elif [[ $1 == "tag" ]]; then 
    tag_release
elif [[ $1 == "build" ]]; then 
    build_push_release_artifacts
elif [[ $1 == "push" ]]; then 
    push_artifacts
elif [[ $1 == "all" ]]; then
    make_tidy_autogen
    write_changelog
    tag_release
    build_push_release_artifacts
    push_artifacts
else
    echo "Usage: $0 [command]"
    echo "Command can be tidy, changelog, tag, build or push."
    echo "Alternatively, 'all' can be specified to do all phases in one."
    echo "To set the version to use, specify the MAJOR, PATCH, and EXTRA environment variables"
fi
