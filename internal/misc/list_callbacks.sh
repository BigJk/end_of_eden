#!/bin/bash

# Parse command line arguments
if [[ $# -ne 1 ]]; then
  echo "This will print all the callbacks with context values from session.go"
  echo "Usage: $0 input_file"
  exit 1
fi

input_file=$1

printf '%-30s: ' "Callback Type"
printf '%s \n' "ctx values"

grep -oE 'Callbacks\[(.+)\]\.Call\(CreateContext\((.+)\)\)' $input_file | while read -r line; do
  # Extract the callback type and context values from the line
  callback_type=$(echo "$line" | sed -E 's/Callbacks\[(.+)\]\.Call\(CreateContext\((.+)\)\)/\1/')
  context_values=$(echo "$line" | sed -E 's/Callbacks\[(.+)\]\.Call\(CreateContext\((.+)\)\)/\2/')

  # Parse the context values into two arrays
  keys=()
  values=()
  while IFS=',' read -ra key_value_pairs; do
    for pair in "${key_value_pairs[@]}"; do
      IFS='=' read -ra kv <<< "$pair"
      key=${kv[0]#[[:space:]]}
      key=${key#\"}
      key=${key%\"}
      keys+=($key)
      values+=(${kv[1]})
    done
  done <<< "$context_values"

  # Print the extracted information in the desired format
  printf '%-30s: ' "$callback_type"
  for (( i=0; i<${#keys[@]}; i += 2 )); do
    printf '%s ' "${keys[i]}"
  done
  printf '\n'
done | sort -k1,1 -t' '
