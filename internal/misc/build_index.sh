#!/bin/bash

# Check if the argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <folder_path>"
  exit 1
fi

folder_path="$1"

# Initialize an empty array
file_index=()

# Function to walk through the directory recursively
walk() {
  local directory="$1"
  for file in "$directory"/*; do
    if [ -d "$file" ]; then
      # If it's a directory, add to the array and recursively call walk function
      file_index+=("  {\"path\": \"$file\", \"isFile\": false}")
      walk "$file"
    elif [ -f "$file" ]; then
      # If it's a file, add to the array
      file_index+=("  {\"path\": \"$file\", \"isFile\": true}")
    fi
  done
}

# Start walking through the directory
walk "$folder_path"

# Create JSON file with the content of the array in the same folder
{
  printf '[\n'
  for ((i = 0; i < ${#file_index[@]}; i++)); do
    printf '%s' "${file_index[i]}"
    if [ $i -ne $(( ${#file_index[@]} - 1 )) ]; then
      printf ','
    fi
    printf '\n'
  done
  printf ']\n'
} > "$folder_path/file_index.json"
