#!/usr/bin/env bash
set -Eeuo pipefail

bindir="${BINDIR:-/usr/local/bin}"
if [ ! -d "$bindir" ]; then
  echo "It seems that the installation target directory $bindir does not exist or is not a directory." 2>&1
  echo "Please make sure the directory $bindir exists." 2>&1
  exit 1
fi
if [ ! -w "$bindir" ]; then
  echo "It seems you do not have a permission to install 'dq' executable file to $bindir." 2>&1
  echo "Please run this script again with appropriate permissions." 2>&1
  exit 1
fi

if \dq --version > /dev/null 2>&1; then
  if [ "$( \command -v dq )" != "$bindir/dq" ] || [ -L "$bindir/dq" ]; then
    echo 'dq is already installed by using another method.' 2>&1
    echo 'Please use the same method if you want to update dq.' 2>&1
    exit 1
  fi
fi

if ! \curl --version > /dev/null 2>&1; then
  echo '"curl" command is required to install dq' 2>&1
  echo 'Please install "curl" command before proceeding.' 2>&1
  exit
fi

get_goos() {
  case "$( uname -s )" in
    "Linux")
      echo "linux"
      ;;
    "Darwin")
      echo "darwin"
      ;;
    "FreeBSD")
      echo "freebsd"
      ;;
    "Windows"*)
      echo "windows"
      ;;
    *)
      echo "unknown system: $( uname -s )" 2>&1
      exit 1
      ;;
  esac
}

get_goarch() {
  case "$( uname -m )" in
    "amd64" | "x86_64")
      echo "amd64"
      ;;
    "i386" | "i686")
      echo "386"
      ;;
    "armv"*)
      echo "arm"
      ;;
    "arm64" | "aarch64")
      echo "arm64"
      ;;
    *)
      echo "unknown architecture: $( uname -m )" 2>&1
      exit 1
      ;;
  esac
}

get_ext_regexp() {
  local goos=$1
  case "$goos" in
    "linux" | "freebsd")
      echo "\.tar\.gz"
      ;;
    "darwin" | "windows")
      echo "\.zip"
      ;;
    *)
      echo "unknown goos: $goos" 2>&1
      exit 1
      ;;
  esac
}

get_ext() {
  local goos=$1
  case "$goos" in
    "linux" | "freebsd")
      echo ".tar.gz"
      ;;
    "darwin" | "windows")
      echo ".zip"
      ;;
    *)
      echo "unknown goos: $goos" 2>&1
      exit 1
      ;;
  esac
}

extract() {
  local path=$1
  local dir=$2
  local ext=$3
  case "$ext" in
    ".tar.gz")
      tar -C "$dir" -xf "$path"
      ;;
    ".zip")
      unzip -q "$path" -d "$dir"
      ;;
    *)
      echo "unknown archive extension: $ext"
      exit 1
      ;;
  esac
}

goos="$( get_goos )"
goarch="$( get_goarch )"
ext_regexp="$( get_ext_regexp "$goos" )"
ext="$( get_ext "$goos" )"

url="$( \curl -fsSL https://api.github.com/repos/bitbears-dev/dq/releases/latest | \
  \grep 'browser_download_url' | \
  \grep "${goos}_${goarch}${ext_regexp}" | \
  cut -d : -f 2-3 | \
  tr -d '"'
)"
shopt -s extglob
url="${url##+( )}" # removes longest matching series of spaces from the front (`shopt -s extglob` is required)
fname=${url##*/}   # removes longest matching series of the pattern from the front (string after the last / will be left)

tmpdir=$(mktemp -d -t dq.XXXXXXXX)
trap 'tear_down' 0

tear_down() {
    : "Clean up tmpdir" && {
        [[ $tmpdir ]] && rm -rf "$tmpdir"
    }
}

echo "Downloading ... "
curl -fSL# --output "$tmpdir/$fname" "$url"
echo "done."

echo -n "Extracting ... "
extract "$tmpdir/$fname" "$tmpdir" "$ext"
echo "done."

echo -n "Installing ... "
mv "$tmpdir/dq" "$bindir"
chmod +x "$bindir/dq"
echo "done."
