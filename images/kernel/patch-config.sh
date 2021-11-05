#!/usr/bin/env bash

OLD_FILE="${1}"
NEW_FILE="${2}"
# The path to where the "patch instructions" are
PATCH_FILE="${3}"

# Note: set_kernel_config and unset_kernel_config are courtesy of @sakaki- on Github
# https://github.com/sakaki-/bcm2711-kernel-bis/blob/master/conform_config.sh.
# Slight modifications have been made to fit this context

set_kernel_config() {
    # flag as $1, value to set as $2, config file as $3
    local TGT="CONFIG_${1#CONFIG_}"
    local REP="${2//\//\\/}"
    local FILE=${3}
    
    if grep -q "${TGT}" ${FILE}; then
        sed "s/^\(${TGT}=.*\|# ${TGT} is not set\)/${TGT}=${REP}/" ${FILE} > ${FILE}.replaced
        mv ${FILE}.replaced ${FILE}
    else
        echo "${TGT}=${REP}" >> ${FILE}
    fi
}

unset_kernel_config() {
    # unsets flag with the value of $1, config file as $2
    local TGT="CONFIG_${1#CONFIG_}"
    local FILE=${3}
    sed "s/^${TGT}=.*/# ${TGT} is not set/" ${NEW_CONFIG_FILE} > ${NEW_CONFIG_FILE}.replaced
    mv ${NEW_CONFIG_FILE}.replaced ${NEW_CONFIG_FILE}
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


# Copy the old config file to the new (overwrite if present), and patch the new one in-place
cp ${OLD_FILE} ${NEW_FILE}
# Add an extra newline to the upstream file if it hasn't got it
# From https://backreference.org/2010/05/23/sanitizing-files-with-no-trailing-newline/
tail -c1 "${NEW_FILE}" | read -r _ || echo >> "${NEW_FILE}"
# Apply patches to the new file
patch_file ${NEW_FILE}

