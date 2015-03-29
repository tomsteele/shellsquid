#!/bin/sh

for OS in "linux" "darwin" "freebsd"; do
    for ARCH in "amd64"; do
        GOOS=$OS  CGO_ENABLED=0 GOARCH=$ARCH godep go build
        FOLDER=shellsquid2.1.0$OS-$ARCH
        ARCHIVE=$FOLDER.tar.gz
        mkdir -p $FOLDER/static
        cp LICENSE $FOLDER
        cp config.json $FOLDER
        if [ $OS = "windows" ] ; then
            cp shellsquid.exe $FOLDER
            rm shellsquid.exe
        else
            cp shellsquid $FOLDER
            rm shellsquid
        fi
            tar -czf $ARCHIVE $FOLDER
        rm -rf $FOLDER
        echo $ARCHIVE
    done
done
