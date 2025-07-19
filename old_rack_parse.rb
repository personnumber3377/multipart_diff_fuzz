require 'rack'
require 'tempfile'
require 'json'

env = {
  "CONTENT_TYPE" => "multipart/form-data; boundary=RubyBoundary",
  "CONTENT_LENGTH" => File.size("oof.bin"),
  "rack.input" => File.open("oof.bin")
}

query_parser = Rack::QueryParser.new(Rack::QueryParser::Params, 32)
parser = Rack::Multipart::Parser.new(
  env["CONTENT_TYPE"].match(/boundary=(.+)/)[1],
  ->(filename, content_type) { Tempfile.new(['RackMultipart', filename]) },  # 🔥 Fix is here
  File.size("oof.bin"),
  query_parser
)

# puts env["CONTENT_TYPE"].match(/boundary=(.+)/)[1]

parser.parse(env["rack.input"])
params = parser.result.params || {}

response = {params: {}, files: {}}
params.each do |k, v|
  if v.respond_to?(:tempfile)
    content = v.tempfile.read
    v.tempfile.rewind
    response[:files][k] = {
      filename: v.original_filename,
      content: content
    }
  else
    response[:params][k] = v
  end
end

puts response.to_json