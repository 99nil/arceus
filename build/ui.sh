#!/bin/bash

set -o errexit
set -u

readonly tmpDir=_tmp
rm -rf $tmpDir && mkdir -p $tmpDir
cd $tmpDir
git clone --depth=1 https://github.com/99nil/arceus-ui.git arceus-ui
cd arceus-ui
echo "Source Downloaded."
yarn install
yarn build
echo "Generated."
cd ../../

mkdir -p static/ui
find static/ui/* | grep -v readme.txt | xargs rm -rf
mv $tmpDir/arceus-ui/build/* static/ui/

echo "Finished."
rm -rf $tmpDir
echo "Cleanup."



