#!/usr/bin/env bash
for image in $(cat preloaded_images)
do
  ctr129 -n moby images import --no-unpack "/tmp/images/${image//[\/\:]/__}.tar"
done
