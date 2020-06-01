#!/usr/bin/env bash
mkdir tmp
for image in $(cat preloaded_images)
do
  docker pull $image
  docker save $image -o "tmp/${image//[\/\:]/__}.tar"
done
