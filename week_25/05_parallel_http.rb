#!/usr/bin/env ruby
require 'net/http'
require 'uri'
require 'json'

# 05. Parallel HTTP Requests in Ruby

puts "=== Ruby Parallel HTTP Requests ==="
puts

urls = [
  'https://jsonplaceholder.typicode.com/posts/1',
  'https://jsonplaceholder.typicode.com/posts/2',
  'https://jsonplaceholder.typicode.com/posts/3',
  'https://jsonplaceholder.typicode.com/posts/4',
  'https://jsonplaceholder.typicode.com/posts/5'
]

# Sequential (slow)
puts "1. Sequential requests:"
start_time = Time.now

urls.each do |url|
  response = Net::HTTP.get_response(URI(url))
  puts "  #{url}: #{response.code}"
end

sequential_time = Time.now - start_time
puts "  Time: #{sequential_time.round(2)}s"
puts

# Parallel (fast)
puts "2. Parallel requests (with threads):"
start_time = Time.now

threads = urls.map do |url|
  Thread.new(url) do |u|
    response = Net::HTTP.get_response(URI(u))
    puts "  #{u}: #{response.code}"
    response
  end
end

responses = threads.map(&:value)  # Wait and collect results
parallel_time = Time.now - start_time

puts "  Time: #{parallel_time.round(2)}s"
puts "  Speedup: #{(sequential_time / parallel_time).round(2)}x"
puts

# With error handling
puts "3. With error handling:"

def fetch_url(url)
  uri = URI(url)
  response = Net::HTTP.get_response(uri)
  
  if response.code == '200'
    data = JSON.parse(response.body)
    { success: true, title: data['title'], url: url }
  else
    { success: false, error: "HTTP #{response.code}", url: url }
  end
rescue => e
  { success: false, error: e.message, url: url }
end

threads = urls.map do |url|
  Thread.new { fetch_url(url) }
end

results = threads.map(&:value)

results.each do |result|
  if result[:success]
    puts "  ✓ #{result[:url]}"
    puts "    Title: #{result[:title][0..50]}..."
  else
    puts "  ✗ #{result[:url]}: #{result[:error]}"
  end
end

puts
puts "✅ Parallel HTTP with error handling complete"

# Key points:
# - Threads good for I/O (HTTP, DB)
# - thread.value waits and returns result
# - Much faster than sequential
# - Handle errors per thread
# - GIL doesn't matter for I/O
