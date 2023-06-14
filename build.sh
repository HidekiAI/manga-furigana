#!/bin/bash
_BUILD_DIR=/tmp/build/manga-furigana
go get -u github.com/ikawaha/kagome/v2/...
go get -u github.com/ikawaha/kagome/v2
go get -u github.com/ikawaha/kagome-dict
go mod tidy

if ! [ -e ${_BUILD_DIR} ]; then
    mkdir -p ${_BUILD_DIR}
fi
cp -r src/* ${_BUILD_DIR}/
cp -r assets/* ${_BUILD_DIR}/
cp go.mod ${_BUILD_DIR}/
if [ -e go.sum ]; then
    cp go.sum ${_BUILD_DIR}/
fi
cd ${_BUILD_DIR}

if ! [ -e /dev/shm/kagome-dict ]; then
    pushd .
    cd /dev/shm
    git clone https://github.com/ikawaha/kagome-dict.git
    popd
fi
if ! [ -e ipa.dict ]; then
    cp /dev/shm/kagome-dict/ipa/ipa.dict . 
fi

# assets (sample) based off of Creative Common License
if ! [ -e ubunchu01_ja ]; then
    mkdir ubunchu01_ja && cd ubunchu01_ja
    unzip ../ubunchu01_ja.zip
fi

go build -o manga_furigana.out nativehost/main.go