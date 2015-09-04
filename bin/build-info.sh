#!/usr/bin/env bash

function generate-build-info() {
	local out=$1
	local snap_stage_name=$2
	local build_number=$3
	local branch=$4
	local commit=$5

	cat > $out <<- EOM
	{
		"stage": "${snap_stage_name}",
		"build_number": "${build_number}",
		"branch": "${branch}",
		"commit": "${commit}",
		"deployed_at": "$(date)"
	}
	EOM
}

generate-build-info $1 $2 $3 $4 $5