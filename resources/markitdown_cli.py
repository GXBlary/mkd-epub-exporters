import sys
import os
import re
from markitdown import MarkItDown

def main():
    if len(sys.argv) < 2:
        print("Usage: markitdown_cli <input_file> [output_file]", file=sys.stderr)
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else None
    
    if not os.path.exists(input_file):
        print(f"Error: Input file '{input_file}' does not exist.", file=sys.stderr)
        sys.exit(1)
        
    try:
        md = MarkItDown()
        result = md.convert(input_file)
        
        # Replace slide number comments with ---
        md_text = result.text_content
        md_text = md_text.replace('\r\n', '\n').replace('\r', '\n')
        md_text = re.sub(r'<!--\s*Slide number:\s*\d+\s*-->', '\n\n---\n\n', md_text)
        while '\n\n\n' in md_text:
            md_text = md_text.replace('\n\n\n', '\n\n')
        
        if output_file:
            # Create directory if it doesn't exist
            out_dir = os.path.dirname(output_file)
            if out_dir and not os.path.exists(out_dir):
                os.makedirs(out_dir, exist_ok=True)
            
            with open(output_file, "w", encoding="utf-8") as f:
                f.write(md_text)
        else:
            # Print to stdout with UTF-8 encoding
            sys.stdout.buffer.write(md_text.encode('utf-8'))
            
    except Exception as e:
        print(f"Error during conversion: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
