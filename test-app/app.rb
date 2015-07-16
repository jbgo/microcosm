require 'rubygems'
require 'bundler/setup'
require 'sinatra'

get "/" do
  "Hello, #{ENV.fetch('WHO', 'You!')}"
end
