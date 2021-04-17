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
mv public/ui public/ui.old && mv $tmpDir/arceus-ui/build public/ui
echo "Finished."
rm -rf $tmpDir public/ui.old
echo "Cleanup."



