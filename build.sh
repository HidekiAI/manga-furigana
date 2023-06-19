#!/bin/bash
# arg1: fully qualified paths of your credentials (i.e. credentials.json) JSON  file paths

_CREDENTIALS_FILE=$1
_PWD=$(pwd)
if [ x"${_CREDENTIALS_FILE}" = x"" ]; then
    echo "Usage: $0 ${_PWD}/credentials.placeholder.json"
    exit 1
fi
set -o nounset  # error when referencing undefined variable

# Even though I'm asking for fully qualified path, I'm going to check for relative path...
_CREDENTIALS_FILE=$(echo ${_CREDENTIALS_FILE} | sed "s|^\.|${_PWD}|g")
if [ ! -e ${_CREDENTIALS_FILE} ]; then
    echo "Error: ${_CREDENTIALS_FILE} does not exist"
    exit 1
fi

# Look for 
# "type": "service_account",
# "universe_domain": "googleapis.com"
_CRED_TYPE_VALIDATION=$(grep "type\":" ${_CREDENTIALS_FILE} | grep "service_account")
_CRED_UD_VALIDATION=$(grep "universe_domain\":" ${_CREDENTIALS_FILE} | grep "googleapis.com")
if [ x"${_CRED_TYPE_VALIDATION}" = x"" ]; then
    echo "Error: ${_CREDENTIALS_FILE} does not contain \"type\": \"service_account\""
    exit 1
fi
if [ x"${_CRED_UD_VALIDATION}" = x"" ]; then
    echo "Error: ${_CREDENTIALS_FILE} does not contain \"universe_domain\": \"googleapis.com\""
    exit 1
fi

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
touch credentials.json
_DIFF=$(diff credentials.json ${_CREDENTIALS_FILE})
if [ x"${_DIFF}" != x"" ]; then
    cp credentials.json credentials.$(date +%Y-%m-%d_%H-%M-%S).bak.json
    cat ${_CREDENTIALS_FILE} > credentials.json
fi
rm credentials.placeholder.json

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

go build -o manga_furigana.out nativehost/*.go

pwd
ls -lAh *.out credentials.*

# now unit-test, note that because the BINARY version as well as Main.go calls the private function init() which expects the
# credentials.json file to exist, we need to create a copy in the current directory as well
echo "# now unit-test"
# copy it for the init() usage
cp credentials.json nativehost/
cp ipa.dict nativehost/
cd nativehost
go test -v ./...
# make sure we only have one copy of credentials.json by removing this temp duplicate...
rm credentials.json
rm ipa.dict
cd .. 