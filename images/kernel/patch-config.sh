#!/bin/bash

# The path to where the "patch instructions" are
PATCH_FILE="./config-patches"

# Note: set_kernel_config and unset_kernel_config are courtesy of @sakaki- on Github
# https://github.com/sakaki-/bcm2711-kernel-bis/blob/master/conform_config.sh.
# Slight modifications have been made to fit this context

set_kernel_config() {
    # flag as $1, value to set as $2, config file as $3
    local TGT="CONFIG_${1#CONFIG_}"
    local REP="${2//\//\\/}"
    local FILE=${3}
    
    if grep -q "${TGT}" ${FILE}; then
        sed -i "s/^\(${TGT}=.*\|# ${TGT} is not set\)/${TGT}=${REP}/" ${FILE}
    else
        echo "${TGT}=${REP}" >> ${FILE}
    fi
}

unset_kernel_config() {
    # unsets flag with the value of $1, config file as $2
    local TGT="CONFIG_${1#CONFIG_}"
    local FILE=${3}
    sed -i "s/^${TGT}=.*/# ${TGT} is not set/" ${NEW_CONFIG_FILE}
}

patch_file() {
    # patches a config file $1 according to the recipe of $PATCH_FILE
    config_file=$1
    echo "Patching ${config_file}..."

    for line in $(cat ${PATCH_FILE} | grep -v "#"); do
        # From https://stackoverflow.com/questions/10638538/split-string-with-bash-with-symbol
        config_key=${line%=*}
        config_value=${line#*=}
        echo "    Applying: ${config_key}=${config_value}"
        if [[ ${config_value} == "n" ]]; then
            unset_kernel_config ${config_key} ${config_file}
        else
            set_kernel_config ${config_key} ${config_value} ${config_file}
        fi
    done
}

for file in ./versioned/*; do
    old_file=${file}
    new_file="./generated/$(basename ${file})"
    

    # Copy the old config file to the new (overwrite if present), and patch the new one in-place
    cp ${old_file} ${new_file}
    # Add an extra newline to the upstream file if it hasn't got it
    # From https://backreference.org/2010/05/23/sanitizing-files-with-no-trailing-newline/
    tail -c1 "${new_file}" | read -r _ || echo >> "${new_file}"
    # Apply patches to the new file
    patch_file ${new_file}
done

