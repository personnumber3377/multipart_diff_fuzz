import email
import sys
import json



# Configure boundary (you can also auto-detect if needed)
boundary = "RubyBoundary"

# Read raw multipart body from stdin
body = sys.stdin.buffer.read()
print(body.decode("ascii"), end="")
# Wrap it in a valid MIME message
wrapped = b"Content-Type: multipart/form-data; boundary=" + boundary.encode() + b"\r\n\r\n" + body

# Parse using the RFC-compliant email parser
msg = email.message_from_bytes(wrapped)

result = {
    "params": {},
    "files": {}
}

for part in msg.walk():
    # print("oof")
    # Skip multipart containers
    if part.get_content_maintype() == "multipart":
        continue
    #     print("part.get_content_maintype(): "+str(part.get_content_maintype()))
    #     print("Invalid...")
    #     exit(1)

    # Get disposition and parameters
    cd = part.get("Content-Disposition")
    # print("part.__repr__(): "+str(part.__repr__()))
    # print("part: "+str(part))
    if not cd or not cd.lower().startswith("form-data"):
        # print("part: "+str(part))
        # print("part.__repr__(): "+str(part.__repr__()))
        # continue  # skip invalid or unsupported parts
        # print("cd: "+str(cd))
        # print("Invalid...")
        exit(1)

    # Parse disposition parameters (e.g., name, filename)
    disposition_params = dict(part.get_params(header="Content-Disposition", failobj=[], unquote=True))
    name = disposition_params.get("name")
    filename = disposition_params.get("filename")

    if name is None:
        # continue  # no usable name -> skip
        # print("feeeeeeeeee")
        # print("Invalid...")
        exit(1)

    content = part.get_payload(decode=True)
    if content is None:
        content = b""

    if filename:
        result["files"][name] = {
            "filename": filename,
            "content": content.decode(errors="replace")
        }
    else:
        result["params"][name] = content.decode(errors="replace")

# Output as JSON
# print(json.dumps(result, indent=2))

