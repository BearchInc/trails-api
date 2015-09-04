#!/usr/bin/env bash

function resolve-dependencies() {
	local gopath=$1
	local gocmd=$2
	local gigo_import_path=github.com/LyricalSecurity/gigo
	local gigo_path=${gopath}/src/${gigo_import_path}

	if [[ ! -d "${gigo_path}" ]]; then
		git clone https://github.com/drborges/gigo ${gigo_path}
		$gocmd get ${gigo_import_path}/...
		$gocmd install ${gigo_import_path}
	fi

	echo "Running gigo from ${gigo_path}"
	brew install wget 2> /dev/null || true
}

resolve-dependencies $1 $2

