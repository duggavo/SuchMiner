mkdir build
cd build

# Download WOWRig (Windows)
mkdir deps_windows
cd deps_windows
wget https://github.com/duggavo/wowrig/releases/download/v6.22.3/xmrig-v6.22.3-win64.zip -O win.zip
unzip -o win.zip
rm win.zip config.json
cd ..

# Download WOWRig (Linux)
mkdir deps_linux
cd deps_linux
wget https://github.com/duggavo/wowrig/releases/download/v6.22.3/xmrig-v6.22.3-lin64.tar.gz -O lin.tar.gz
tar -xzf lin.tar.gz
rm lin.tar.gz config.json
cd ..

# Build for Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" ..
zip suchminer-windows-amd64.zip suchminer.exe deps_windows/*

# Build for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" ..
tar -cJf suchminer-linux-amd64.tar.xz suchminer deps_linux/*
