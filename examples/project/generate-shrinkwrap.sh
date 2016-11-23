# Generate npm-shrinkwrap.json with different versions of npm,
# to compare differences in file format.
#
# Example usage: ./generate-shrinkwrap.sh 2.11.3
#
# This will first install npm version 2.11.3 (globally),
# and then generate the shrinkwrap file.

VER=$1

npm install -g npm@$VER
rm -rf node_modules
rm -rf npm-shrinkwrap.json
npm install
npm shrinkwrap
mv npm-shrinkwrap.json npm-shrinkwrap-${VER}.json
