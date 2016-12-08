#!/bin/sh

make update

nginx &

while true; do date; make update; sleep 20; done
