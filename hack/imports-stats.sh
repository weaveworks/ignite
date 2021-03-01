#!/usr/bin/env bash

BINARIES=${@:-"ignite ignited ignite-spawn"}
PRINT=${PRINT:-"false"}

for binary in ${BINARIES}; do
    echo "${binary} imports: $(go list -f '{{ join .Deps "\n" }}' ./cmd/${binary} | grep "\." | wc -l)"
    echo "${binary} size: $(ls -l bin/${binary} | awk '{print $5}') ($(ls -lh bin/${binary} | awk '{print $5}'))"
    echo "----"
    if [[ ${PRINT} == "true" ]]; then
        echo "A list of ${binary} imports"
        go list -f '{{ join .Deps "\n" }}' ./cmd/${binary} | grep "\."
        echo "----"
    fi
done
