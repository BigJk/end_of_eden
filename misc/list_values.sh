#!/bin/bash

# Parse command line arguments
if [[ $# -ne 1 ]]; then
  echo "This will print all the functions and values from lua.go"
  echo "Usage: $0 input_file"
  exit 1
fi

input_file=$1

printf "%-30s : type\n" "name"
while IFS= read -r line; do
  if [[ $line =~ l.SetGlobal\(\"(.*)\",\ *(.*)\) ]]; then
    key=${BASH_REMATCH[1]}
    value=${BASH_REMATCH[2]}
    if [[ $value == l.NewFunction* ]]; then
      printf "%-30s : function\n" "$key"
    else
      printf "%-30s : value\n" "$key"
    fi
  fi
done < "$input_file"