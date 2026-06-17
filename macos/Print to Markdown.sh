#!/bin/bash
# macOS PDF Service: Print to Markdown
# Receives the generated PDF path as the first argument ($1)

PDF_PATH="$1"

# 1. Ask the user for the output path using native AppleScript
OUT_PATH=$(osascript -e 'tell application "System Events" to POSIX path of (choose file name with prompt "Save Markdown file as:" default name "document.md")' 2>/dev/null)

# 2. Exit if the user cancelled the dialog
if [ -z "$OUT_PATH" ]; then
    exit 0
fi

# 3. Create a temporary folder to run the converter
TEMP_DIR=$(mktemp -d)

# 4. Execute the converter CLI
/usr/local/bin/markitdown-cli -out "$TEMP_DIR" -format md "$PDF_PATH"

# 5. Move the converted file to the user's selected path
CONV_FILE=$(find "$TEMP_DIR" -type f \( -name "*.md" \) | head -n 1)
if [ -f "$CONV_FILE" ]; then
    mv "$CONV_FILE" "$OUT_PATH"
fi

# 6. Cleanup temp folder
rm -rf "$TEMP_DIR"
