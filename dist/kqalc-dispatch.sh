#!/bin/sh
# Picks the binary matching the host architecture.
# krunner-plugininstaller writes this script's path into the D-Bus service
# file, so D-Bus activation lands here rather than on a per-arch binary.
set -e

dir=$(dirname "$(readlink -f "$0")")

case "$(uname -m)" in
	x86_64)        exec "$dir/kqalc-amd64" "$@" ;;
	aarch64|arm64) exec "$dir/kqalc-arm64" "$@" ;;
	*)             echo "kqalc: unsupported architecture $(uname -m)" >&2; exit 1 ;;
esac
