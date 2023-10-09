#!/bin/bash

set -e

if
    command -v pho >/dev/null
then
    echo "Pho is already installed! Exiting..."
    exit 1
fi

install_dir="${HOME}/.local/bin"
custom_install_dir=false
if [ -n "$1" ]; then
    install_dir=$1
    custom_install_dir=true
fi
install_dir="${install_dir%/}"
if ! [[ -d "${install_dir}" ]]; then
    echo "[error] Installation directory '${install_dir}' does not exist."
    exit 1
fi
if [[ $custom_install_dir ]]; then
    echo "Kindly ensure that the installation directory exists in \$PATH variable"
fi

install_path="${install_dir}/pho"
echo "Installation path set as '${install_dir}'"

arch=$(uname -m)
case "${arch}" in
"x86_64")
    arch="amd64"
    ;;
"i386")
    arch="386"
    ;;
"i686")
    arch="386"
    ;;
"armhf")
    arch="arm64"
    ;;
"aarch64")
    arch="arm64"
    ;;
esac
echo "Architecture found to be '${arch}'"

release_url="https://api.github.com/repos/zyrouge/pho/releases/latest"
download_url=$(
    curl -Ls --fail "${release_url}" |
        grep -E "\"browser_download_url\".*pho-${arch}\"" |
        sed -nr 's/.*"([^"]+)"$/\1/p'
)
if [[ "${download_url}" == "" ]]; then
    echo "[error] Unsupported platform or architecture."
    exit 1
fi

echo "Downloading binary from '${download_url}'..."
curl --fail -Ls -o "${install_path}" "${download_url}"
chmod +x "${install_path}"

echo "Installation succeeded!"
echo "You can get started by using 'pho init' to initialize Pho!"
