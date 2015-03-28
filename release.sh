#!/bin/sh

cd client
grunt build
cd ..

for OS in "linux" "darwin" "freebsd"; do
    for ARCH in "amd64"; do
        GOOS=$OS  CGO_ENABLED=0 GOARCH=$ARCH go build
        FOLDER=shellsquid2.0.2$OS-$ARCH
        ARCHIVE=$FOLDER.tar.gz
        mkdir -p $FOLDER/client/dist
        cp LICENSE $FOLDER
        cp config.json $FOLDER
        cp -R client/dist/* $FOLDER/client/dist
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
