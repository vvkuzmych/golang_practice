#!/usr/bin/env ruby

# 07. Race Condition Demo in Ruby

puts "=== Ruby Race Condition Demo ==="
puts

# Problem: Race condition
puts "1. Race condition (BROKEN CODE):"
counter = 0

threads = 100.times.map do
  Thread.new do
    100.times { counter += 1 }
  end
end

threads.each(&:join)

puts "  Expected: 10000"
puts "  Got:      #{counter}"
puts "  Lost:     #{10000 - counter} increments!"
puts

# Why it happens:
puts "2. Why race condition happens:"
puts "  counter += 1 is actually 3 operations:"
puts "    1. Read current value"
puts "    2. Add 1"
puts "    3. Write new value"
puts
puts "  Thread A:        Thread B:"
puts "  Read (0)         Read (0)"
puts "  Add 1            Add 1"
puts "  Write (1)        Write (1)  ← Lost increment!"
puts

# Solution: Mutex
puts "3. Solution with Mutex:"
mutex = Mutex.new
counter = 0

threads = 100.times.map do
  Thread.new do
    100.times do
      mutex.synchronize do
        counter += 1
      end
    end
  end
end

threads.each(&:join)

puts "  Expected: 10000"
puts "  Got:      #{counter}"
puts "  Lost:     0 (protected by mutex!)"
puts

# More complex race: Array
puts "4. Race condition with Array:"
array = []

threads = 10.times.map do |i|
  Thread.new do
    100.times { array << i }
  end
end

threads.each(&:join)

puts "  Expected length: 1000"
puts "  Got:             #{array.length}"
puts "  Lost:            #{1000 - array.length} items (maybe)"
puts

# Solution: Thread-safe data structures
puts "5. Solution: Mutex-protected array"
mutex = Mutex.new
array = []

threads = 10.times.map do |i|
  Thread.new do
    100.times do
      mutex.synchronize { array << i }
    end
  end
end

threads.each(&:join)

puts "  Expected length: 1000"
puts "  Got:             #{array.length}"
puts "  Perfect! ✅"
puts

puts "✅ Race condition demo complete"
puts
puts "💡 Always use Mutex for shared mutable state!"

# Key points:
# - Race conditions happen with shared state
# - counter += 1 is NOT atomic
# - GIL doesn't prevent race conditions!
# - Use Mutex to protect critical sections
# - Thread-safe: Queue, SizedQueue
# - Not thread-safe: Array, Hash, variables
