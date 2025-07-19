require 'rack'
require 'tempfile'
require 'json'
require 'stringio'

# Read from STDIN
input_data = STDIN.read

env = {
  "CONTENT_TYPE" => "multipart/form-data; boundary=RubyBoundary",
  "CONTENT_LENGTH" => input_data.bytesize,
  "rack.input" => StringIO.new(input_data)
}

query_parser = Rack::QueryParser.new(Rack::QueryParser::Params, 32)

parser = Rack::Multipart::Parser.new(
  env["CONTENT_TYPE"].match(/boundary=(.+)/)[1],
  ->(filename, content_type) { Tempfile.new(['RackMultipart', filename]) },
  input_data.bytesize,
  query_parser
)

parser.parse(env["rack.input"])
params = parser.result.params || {}

response = { params: {}, files: {} }

params.each do |key, value|
  if value.is_a?(Hash) &&
     value.key?(:filename) &&
     value.key?(:tempfile) &&
     value[:tempfile].respond_to?(:read)

    content = value[:tempfile].read
    value[:tempfile].rewind

    response[:files][key] = {
      filename: value[:filename],
      content: content
    }
  else
    response[:params][key] = value
  end
end

puts JSON.pretty_generate(response)