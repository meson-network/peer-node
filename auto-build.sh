#!/bin/sh
# Copyright 2020 Daqnext Foundation Ltd.

VERSION="v3.0.0"
COPY_FILES=("configs/pro.json")

generate_tar() {
  mkdir "./build/$1/configs"
  cp  "configs/pro.json" "./build/$1/configs/pro.json" && cd build && tar -czvf "$1.tar.gz" $1 && rm -rf $1 && cd ..|| exit
}

generate_zip(){
  mkdir "./build/$1/configs"
  cp "configs/pro.json" "./build/$1/configs/pro.json" && cd build && zip -r "$1.zip" $1 && rm -rf $1 && cd ..|| exit
}

rm -f -R ./build
mkdir build

#todo mac is for dev, not for pro
echo "Compiling MAC     x86_64 version"
DIR="meson-darwin-amd64" && GOOS=darwin GOARCH=amd64 go build -o "./build/${DIR}/meson" && generate_tar ${DIR}

#echo "Compiling Windows x86_64 version"
#DIR="meson-windows-386" && GOOS=windows GOARCH=386   go build -o "./build/${DIR}/meson.exe" && generate_zip ${DIR}
#DIR="meson-windows-amd64" && GOOS=windows GOARCH=amd64 go build -o "./build/${DIR}/meson.exe" && generate_zip ${DIR}

#echo "Compiling Linux   x86_64 version"
#DIR="meson-linux-386"   &&  GOOS=linux GOARCH=386   go build -o "./build/${DIR}/meson" && generate_tar ${DIR}
#DIR="meson-linux-amd64" &&  GOOS=linux GOARCH=amd64 go build -o "./build/${DIR}/meson" && generate_tar ${DIR}
#
#echo "Compiling ARM64    version"
#DIR="meson-linux-arm64" &&  GOOS=linux GOARCH=arm64 go build -o "./build/${DIR}/meson" && generate_tar ${DIR}