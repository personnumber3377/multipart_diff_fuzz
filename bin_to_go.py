
import sys

def write_escaped_hex(input_path, output_path):
    with open(input_path, "rb") as f:
        binary_data = f.read()

    with open(output_path, "w") as f:
        escaped = ''.join(f'\\x{byte:02x}' for byte in binary_data)
        start = '''go test fuzz v1\n[]byte("'''
        end = "\")"
        f.write(start+escaped+end)

# Example usage:
if len(sys.argv) != 3:
    print("Usage: "+str(sys.argv[0]), " INFILE OUTFILE")
    exit(1)
write_escaped_hex(sys.argv[1], sys.argv[2])
