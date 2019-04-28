#!/usr/bin/env ruby

# require 'net/http'
require 'httparty'
require 'json'

file = File.open('BL/porn/domains')

i = 0
while line = file.gets
#   puts line
  begin
    response = HTTParty.post('http://localhost:8080/insert/', body: { host: line.chomp})
  rescue StandardError => e
    puts('Error handling request')
    gets
  end
#   puts response.code
#   puts response.body
#   gets
    if i % 10000 == 0
        puts i
    end
    i = i+1

   if response.code != 200
    puts "Error handling #{line}"
    gets
   end
end
