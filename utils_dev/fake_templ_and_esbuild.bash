#ls
#echo "/*" >> tmp/fake.templ
#date +%F" "%T%n"*/" >> tmp/fake.templ
cd utils_dev

OPEN="/*"
CLOSE="*/"
DATE=$(date +%F" "%T)
echo $OPEN$DATE$CLOSE >> fake.templ
#echo "add to fake templ"

# build go binary, valja povremeno da se uradi radi testiranja
#go build -o bin src/main.go

# npx esbuild ../src/react/*.jsx --outdir=../assets/assignments/ --minify --bundle --platform=node --global-name=bundle
npx esbuild ../src/ext/site/*.* --outdir=../assets/ --minify
# npx esbuild ../src/ext/react/*.* --outdir=../assets/assignments/ --minify --bundle --global-name=bundle
cd ..