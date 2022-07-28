#!/bin/sh
# Copyright 2020 Daqnext Foundation Ltd.

generate_tar() {
  mkdir "./build/$1/configs"
  cp  "configs/default.toml" "./build/$1/configs/default.toml" && cd build && tar -czvf "$1.tar.gz" $1 && rm -rf $1 && cd ..|| exit
}

generate_zip(){
  mkdir "./build/$1/configs"
  cp "configs/default.toml" "./build/$1/configs/default.toml" && cd build && zip -r "$1.zip" $1 && rm -rf $1 && cd ..|| exit
}

rm -f -R ./build
mkdir build

#todo mac is for dev, not for pro
#echo "Compiling MAC x86_64 version"
#DIR="meson_cdn-darwin-amd64" && GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o "./build/${DIR}/meson_cdn" && cp "daemon-darwin-amd64" "./build/${DIR}/service" && generate_tar ${DIR}

echo "Compiling Linux bit 64 version"
DIR="meson_cdn-linux-amd64" &&  GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build -a -o "./build/${DIR}/meson_cdn" && cp "daemon-linux-amd64" "./build/${DIR}/service" && generate_tar ${DIR}

echo "Compiling ARM64 version"
DIR="meson_cdn-linux-arm64" &&  GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-musl-gcc CGO_LDFLAGS="-static"  go build -a -o "./build/${DIR}/meson_cdn" && cp "daemon-linux-arm64" "./build/${DIR}/service" && generate_tar ${DIR}

echo "Compiling ARM32 version"
DIR="meson_cdn-linux-arm" &&  GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-musleabihf-gcc CGO_LDFLAGS="-static"  go build -a -o "./build/${DIR}/meson_cdn" && cp "daemon-linux-arm32" "./build/${DIR}/service" && generate_tar ${DIR}
