#!/usr/bin/env bash

SCRIPT_DIR=$( dirname "${BASH_SOURCE[0]}" )
cd ${SCRIPT_DIR}/..

Resources="VM Image Kernel"
for Resource in ${Resources}; do
    resource=$(echo "${Resource}" | awk '{print tolower($0)}')
    sed -e "s|Resource|${Resource}|g;s|resource|${resource}|g;/build ignore/d" \
        pkg/client/client_resource_template.go > \
        pkg/client/zz_generated.client_${resource}.go
done
