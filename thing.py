

data = "--RubyBoundary\r\nContent-Disposition: form-data; name=\"username\"\r\n\r\nalice\r\n--RubyBoundary\r\nContent-Disposition: form-data; name=\"file\"; filename=\"hello.txt\"\r\nContent-Type: text/plain\r\n\r\nhello world\r\n--RubyBoundary--\r\n"
print(data, end="")
