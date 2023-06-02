#!/usr/bin/bash

mkdir -p ~/tmp

wget https://go.dev/dl/go1.19.7.linux-amd64.tar.gz -P ~/tmp

tar -xf ~/tmp/go1.19.7.linux-amd64.tar.gz -C ~/tmp

rm -rf ~/tmp/go1.19.7.linux-amd64.tar.gz
