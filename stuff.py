data = '''--RubyBoundary
Content-Disposition: form-data; name="file"; filename="hello.txt"
Content-Type: text plain

hello world
--RubyBoundary--
'''


data = data.replace("\n", "\r\n")

print(data, end="")



