# Developer Docs

## Pre-Commit Hook

```bash
#!/bin/bash

# Format lua
./format-lua.sh

# Generate lua docs
./update-docs.sh

# Add the modified files to the commit
git add .

# Exit with success
exit 0
```