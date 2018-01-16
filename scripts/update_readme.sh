#!/bin/sh
L0_VERSION=$(git describe --tags)

# Remove v
L0_VERSION=`echo "$L0_VERSION" | sed "s/v//"`

echo Updating to $L0_VERSION
sed -i.bac 's/0\.10\../'$L0_VERSION'/g' ../README.md
rm ../README.md.bac