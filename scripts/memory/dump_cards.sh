#!/bin/bash

# This file help zip cards (images) into seperate archives. Help to make cards dump.
# Only for low-res image in memory-server.
# For storage-server (high-res) will be another file.
cd /var/cards
lastFolder=$(ls -tr | tail -1)
echo "folder $lastFolder ignored"

for x in `ls -tr -I$lastFolder`; do
  filename="/root/dumps/cards/$x.tar.gz"
  echo "working with folder $x and file $filename"
  if [ -e "$filename" ]; then
      echo "archive $filename already exists. skip"
  else
      echo "starting collecting part $x into $filename"
      tar -zcvf $filename $x
      echo "successfully zipped $x into $filename"
      echo "success $x"
  fi
done
