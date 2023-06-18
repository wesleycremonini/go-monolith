#!/bin/bash

cd public
./tailwindcss -i styles.css -o global.css --minify
# ./tailwindcss -i styles.css -o global.css --watch
cd ..
go run .