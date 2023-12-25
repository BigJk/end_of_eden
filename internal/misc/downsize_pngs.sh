#!/bin/bash

input_dir="$1"

if [ ! -d "$input_dir" ]; then
    echo "Usage: $0 <input_directory>"
    exit 1
fi

# Loop through all PNG files in the specified directory
for file in "$input_dir"/*.png; do
    if [ -f "$file" ]; then
        filename=$(basename -- "$file")
        filename_no_extension="${filename%.*}"
        output_file="${file%.png}.jpg"

        # Convert PNG to JPG with quality 30 and save in the same directory
        convert "$file" -quality 30 "$output_file"

        echo "Converted $file to $output_file"
    fi
done