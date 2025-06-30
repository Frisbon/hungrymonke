#!/usr/bin/env bash

# Usage: ./open-npm.sh [path/to/webui]
# If no path is supplied the script will attempt to locate the
# repository root via git and use the "webui" directory inside it.

set -euo pipefail

if [ "$#" -gt 0 ]; then
  WEBUI_PATH="$(realpath "$1")"
else
  if git_root=$(git rev-parse --show-toplevel 2>/dev/null); then
    WEBUI_PATH="$git_root/webui"
  else
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    WEBUI_PATH="$SCRIPT_DIR/webui"
  fi
fi

echo "Mounting path: $WEBUI_PATH"

if [ ! -d "$WEBUI_PATH" ]; then
  echo "Error: webui directory does not exist at $WEBUI_PATH"
  exit 1
fi

# Run the Docker container inside the webui directory
docker run --rm -it \
  -v "$WEBUI_PATH:/app" \
  -p 8081:8081 \
  -w /app \
  node:lts bash
