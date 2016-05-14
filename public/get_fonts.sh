#!/bin/bash

for url in $(grep -P -o 'url\(.*?\)' $1 | grep -Po 'http.*2')
do
    cd ./fonts && wget $url && cd ..
done

