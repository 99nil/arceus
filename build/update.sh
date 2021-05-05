#!/bin/bash

tmpDir=_tmp
rm -rf $tmpDir && mkdir -p $tmpDir
cd $tmpDir || exit
git clone --depth=1 https://github.com/zc2638/arceus-ui.git arceus-ui
cd arceus-ui || exit
echo "Source Downloaded."
yarn install || exit
yarn build || exit
echo "Generated."
cd ../../

if [ -d "static/ui" ]; then
  mv static/ui static/ui.old
fi
mv $tmpDir/arceus-ui/build static/ui

echo "Finished."
rm -rf $tmpDir static/ui.old
echo "Cleanup."



