#!/usr/bin/env ruby

# 01. Basic Threads in Ruby

puts "=== Ruby Threads - Basic ==="
puts

# 1. Simple thread
puts "1. Simple thread:"
thread = Thread.new do
  puts "  Hello from thread!"
  sleep 0.5
  puts "  Thread finished"
end

puts "  Main thread continues..."
thread.join  # Wait for thread to finish
puts

# 2. Thread with parameters
puts "2. Thread with parameters:"
thread = Thread.new(5, "test") do |num, text|
  puts "  Received: num=#{num}, text=#{text}"
end
thread.join
puts

# 3. Multiple threads
puts "3. Multiple threads:"
threads = []
5.times do |i|
  threads << Thread.new(i) do |n|
    sleep(rand * 0.5)
    puts "  Thread #{n} finished"
  end
end

threads.each(&:join)
puts

# 4. Thread return value
puts "4. Thread return value:"
thread = Thread.new do
  sleep 0.2
  42  # Return value
end

result = thread.value  # Waits and returns value
puts "  Result: #{result}"
puts

# 5. Thread status
puts "5. Thread status:"
t = Thread.new { sleep 10 }
puts "  Alive? #{t.alive?}"
puts "  Status: #{t.status}"
t.kill
puts "  After kill - Alive? #{t.alive?}"
puts

# 6. Thread-local variables
puts "6. Thread-local variables:"
Thread.current[:user_id] = 100

thread = Thread.new do
  Thread.current[:user_id] = 200
  puts "  Inside thread: user_id = #{Thread.current[:user_id]}"
end

thread.join
puts "  Main thread: user_id = #{Thread.current[:user_id]}"
puts

puts "✅ Ruby threads complete"

# Key points:
# - OS threads (heavier than goroutines)
# - GIL limits CPU parallelism
# - Good for I/O operations
# - Thread.new { block }
# - thread.join to wait
# - thread.value to get result
