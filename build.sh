#!/bin/sh
mkdir Bundle

export CGO_ENABLED=1;
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
go build -o ./Bundle/bin64.exe

export GOOS=windows
export GOARCH=386
export CC=i686-w64-mingw32-gcc
export CXX=i686_64-w64-mingw32-g++
go build -o ./Bundle/bin32.exe

cp -r ./Website ./Bundle/Website
