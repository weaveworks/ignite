#!/usr/bin/env bash
SCRIPT_DIR=$( dirname "${BASH_SOURCE[0]}" )
cd ${SCRIPT_DIR}/..

cat > graph.dot <<EOD
digraph {
	graph [overlap=false, size=14];
	root="$(go list -m)";
	node [  shape = plaintext, fontname = "Helvetica", fontsize=24];
	"$(go list -m)" [style = filled, fillcolor = "#E94762"];
    $(go mod graph | sed -Ee 's/@[^[:blank:]]+//g' | sort | uniq | awk '{print "\""$1"\" -> \""$2"\""};')
}
EOD

sed -i 's+\("github.com/[^/]*/\)\([^"]*"\)+\1\\n\2+g' graph.dot
sfdp -Tsvg -o docs/dependencies.svg graph.dot
rm graph.dot
