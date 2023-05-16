#!/bin/bash

for file in ./assets/scripts/**/*.lua; do
  echo "Formatting: $file"
  lua-format "$file" -i
done