#!/bin/bash
set -x

echo "Windows"
GOOS=windows GOARCH=amd64 GO_EXTLINK_ENABLED=0 CGO_ENABLED=1 go build
if [[ $? -eq 0 ]]
then
  compress -f pcap2ts.exe
fi

echo "linux"
GOOS=linux GOARCH=amd64 GO_EXTLINK_ENABLED=0 CGO_ENABLED=1 go build  
if [[ $? -eq 0 ]]
then
  mv pcap2ts pcap2ts_linux
  compress -f pcap2ts_linux
fi

echo "mac"
GOOS=darwin GOARCH=amd64 GO_EXTLINK_ENABLED=0 CGO_ENABLED=1 go build 
if [[ $? -eq 0 ]]
then
  mv pcap2ts pcap2ts_mac
  compress -f pcap2ts_mac
fi


