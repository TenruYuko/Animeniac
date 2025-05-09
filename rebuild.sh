#!/bin/bash
echo "Stopping Seanime..."
pkill -9 -f seanime || true

echo "Removing cache..."
rm -rf /aeternae/functional/dockers/animechanical/data/cache/*
rm -rf /aeternae/functional/dockers/animechanical/seanime-web/.next

# Navigate to the web directory
cd /aeternae/functional/dockers/animechanical/seanime-web

echo "Installing dependencies..."
npm install

echo "Building the web interface..."
npm run build

echo "Removing old web files..."
rm -rf /aeternae/functional/dockers/animechanical/web/*

echo "Copying new web files..."
mkdir -p /aeternae/functional/dockers/animechanical/web
cp -r /aeternae/functional/dockers/animechanical/seanime-web/out/* /aeternae/functional/dockers/animechanical/web/

echo "Web UI rebuild complete"

# Navigate back to project root
cd /aeternae/functional/dockers/animechanical
