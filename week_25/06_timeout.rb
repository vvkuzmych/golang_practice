#!/usr/bin/env ruby
require 'timeout'

# 06. Timeout Pattern in Ruby

puts "=== Ruby Timeout ==="
puts

# 1. Basic timeout
puts "1. Basic timeout:"
begin
  Timeout.timeout(2) do
    puts "  Starting long operation..."
    sleep 1
    puts "  ✓ Completed in time"
  end
rescue Timeout::Error
  puts "  ✗ Operation timed out!"
end
puts

# 2. Timeout that exceeds
puts "2. Timeout that exceeds:"
begin
  Timeout.timeout(1) do
    puts "  Starting long operation..."
    sleep 3
    puts "  This won't print"
  end
rescue Timeout::Error => e
  puts "  ✗ Timed out after 1 second!"
end
puts

# 3. Thread with timeout
puts "3. Thread with timeout:"
thread = Thread.new do
  sleep 5
  "result"
end

# Wait with timeout
result = nil
timeout_thread = Thread.new do
  sleep 1  # 1 second timeout
end

if thread.join(1)  # Wait max 1 second
  result = thread.value
  puts "  ✓ Thread completed: #{result}"
else
  thread.kill
  puts "  ✗ Thread timed out and killed"
end
puts

# 4. Custom timeout with Thread
puts "4. Custom timeout implementation:"

def with_timeout(seconds)
  result = nil
  error = nil
  
  thread = Thread.new do
    begin
      result = yield
    rescue => e
      error = e
    end
  end
  
  if thread.join(seconds)
    raise error if error
    result
  else
    thread.kill
    raise Timeout::Error, "Operation timed out after #{seconds}s"
  end
end

# Success case
begin
  result = with_timeout(2) do
    sleep 0.5
    "Success!"
  end
  puts "  ✓ Result: #{result}"
rescue Timeout::Error => e
  puts "  ✗ #{e.message}"
end

# Timeout case
begin
  result = with_timeout(1) do
    sleep 3
    "Won't complete"
  end
  puts "  Result: #{result}"
rescue Timeout::Error => e
  puts "  ✗ #{e.message}"
end

puts
puts "✅ Ruby timeout complete"

# Key points:
# - Timeout.timeout(seconds) { block }
# - Raises Timeout::Error
# - thread.join(timeout) for manual control
# - thread.kill to force stop (dangerous!)
# - Use for HTTP, DB, long operations
