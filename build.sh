#!/usr/bin/env bash

app="gitlog"
deployment="./deploy"
tmp="tmp"
zip="zip"

echo "building gitlog for deployment"
declare -a systems=("darwin" "linux" "windows")

rm -rf "$deployment"
mkdir -p "$deployment/$tmp"
mkdir -p "$deployment/$zip"

for s in "${systems[@]}"
do
    echo "building $s..."
    cmd="$app"
    if [ "$s" == "windows" ]; then
        cmd="$app.exe"
    fi
    CGO_ENABLED=0 GOARCH=amd64 GOOS="$s" packr2 build -o "$deployment/$tmp/$cmd" "./"
    target=$(ls "$deployment/$tmp/")
    zip -j "$deployment/zip/gitlog.$s.zip" "$deployment/$tmp/$target"
    rm -f "$deployment/$tmp/$cmd"
done

