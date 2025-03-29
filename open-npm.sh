#!/bin/bash

# Hardcode the absolute path to the webui directory in Windows format
WEBUI_PATH="C:/Users/sasha/Documents/GitHub/hungrymonke/webui"

# Debug: Print the path being used
echo "Mounting path: $WEBUI_PATH"

# Verify the path exists on the host
if [ -d "$WEBUI_PATH" ]; then
  echo "webui directory exists on host"
else
  echo "Error: webui directory does not exist at $WEBUI_PATH"
  exit 1
fi

# Run the Docker container without the -w flag
docker run --rm -it \
  -v "$WEBUI_PATH:/app" \
  -p 8081:8081 \
  node:lts \
  bash -c "ls -la /app && echo 'Current directory: $PWD' && bash"