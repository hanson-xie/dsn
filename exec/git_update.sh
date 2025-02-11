#!/bin/bash

# Navigate to the script directory
cd "$(dirname "$0")/.." || exit 1

# Ensure we are in a git repository
if [ ! -d .git ]; then
    echo "Error: This is not a git repository!"
    exit 1
fi

# Pull the latest changes from the current branch
echo "Fetching the latest changes..."
git fetch --all

echo "Pulling updates from the current branch..."
git pull

# Check for errors
if [ $? -ne 0 ]; then
    echo "Error: Failed to update the repository."
    exit 1
fi

echo "Git repository successfully updated."