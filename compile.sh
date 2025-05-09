# Navigate to the web directory
cd /aeternae/functional/dockers/animechanical/seanime-web

# Install dependencies
npm install

# Build the web interface
npm run build

# Create the web directory at project root (if it doesn't exist)
mkdir -p /aeternae/functional/dockers/animechanical/web

# Move the built files to the web directory
cp -r /aeternae/functional/dockers/animechanical/seanime-web/out/* /aeternae/functional/dockers/animechanical/web/

# Navigate back to project root
cd /aeternae/functional/dockers/animechanical

# Build the server based on your platform
# For Linux:
go build -o seanime -trimpath -ldflags="-s -w"