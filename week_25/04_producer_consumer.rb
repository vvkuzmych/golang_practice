#!/usr/bin/env ruby
require 'thread'

# 04. Producer/Consumer Pattern in Ruby

puts "=== Ruby Producer/Consumer ==="
puts

queue = Queue.new
mutex = Mutex.new
counter = 0

# Producer thread
producer = Thread.new do
  10.times do |i|
    sleep 0.1
    queue << i
    puts "Producer: added #{i}"
  end
  queue << :done
  puts "Producer: finished"
end

# Consumer threads
consumers = 3.times.map do |id|
  Thread.new do
    loop do
      item = queue.pop
      break if item == :done
      
      puts "  Consumer #{id}: processing #{item}"
      sleep 0.2
      
      mutex.synchronize { counter += 1 }
    end
    
    # Put :done back for other consumers
    queue << :done
    puts "  Consumer #{id}: shutting down"
  end
end

# Wait for completion
producer.join
consumers.each(&:join)

puts
puts "✅ Processed #{counter} items"
puts

# Advanced: Multiple producers, multiple consumers
puts "=== Multiple Producers & Consumers ==="
puts

queue = Queue.new
items_produced = 0
items_consumed = 0
prod_mutex = Mutex.new
cons_mutex = Mutex.new

# Multiple producers
producers = 3.times.map do |id|
  Thread.new do
    5.times do |i|
      item = "P#{id}-#{i}"
      queue << item
      prod_mutex.synchronize { items_produced += 1 }
      puts "Producer #{id}: #{item}"
      sleep 0.05
    end
  end
end

# Multiple consumers
consumers = 2.times.map do |id|
  Thread.new do
    loop do
      item = queue.pop(true) rescue break  # Non-blocking pop
      puts "  Consumer #{id}: #{item}"
      cons_mutex.synchronize { items_consumed += 1 }
      sleep 0.1
    end
  end
end

# Wait for producers
producers.each(&:join)

# Wait a bit for consumers to finish
sleep 1

puts
puts "Produced: #{items_produced}, Consumed: #{items_consumed}"
puts
puts "✅ Multiple producers/consumers complete"

# Key points:
# - Queue is thread-safe
# - queue.pop blocks until item available
# - queue << item to add
# - Use sentinel value (:done) to stop
# - Mutex for shared counter
