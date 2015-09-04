#!/usr/bin/env bash

function download-appengine() {
	local cache_dir=$1; [[ $cache_dir != /* ]] && cache_dir=$(pwd)/${cache_dir}
	# For more on getting the latest version of appengine SDK: http://bekt.github.io/p/gae-sdk/#sthash.v89zGu1k.LHxqfoBi.dpuf
	local appengine_version=$(curl -s https://appengine.google.com/api/updatecheck | awk -F '\"' '/release/ {print $2}')
	local appengine_path=${cache_dir}/go_appengine
	local appengine_sdk=go_appengine_sdk_linux_386-${appengine_version}.zip; [[ $(uname) == "Darwin" ]] && appengine_sdk=go_appengine_sdk_darwin_386-${appengine_version}.zip

	if [[ ! -d "${appengine_path}" ]]; then
		echo "Downloading appengine ${appengine_version} at ${appengine_path}"

		wget -nv -P ${cache_dir} "https://storage.googleapis.com/appengine-sdks/featured/${appengine_sdk}" --no-check-certificate
		unzip -qq -d ${cache_dir} ${cache_dir}/${appengine_sdk}
		exit 0
	fi

	echo "Running appengine ${appengine_version} from ${appengine_path}"
}

download-appengine $1
