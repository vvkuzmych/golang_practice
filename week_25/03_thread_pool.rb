#!/usr/bin/env ruby
require 'thread'

# 03. Thread Pool in Ruby

puts "=== Ruby Thread Pool ==="
puts

WORKER_COUNT = 3
queue = Queue.new

# Create worker threads
workers = WORKER_COUNT.times.map do |i|
  Thread.new do
    loop do
      task = queue.pop
      break if task == :stop
      
      puts "Worker #{i}: processing task #{task}"
      sleep 0.2  # Simulate work
      puts "Worker #{i}: task #{task} done"
    end
    puts "Worker #{i}: shutting down"
  end
end

# Add tasks to queue
puts "Adding 10 tasks..."
10.times { |i| queue << i }
puts

# Stop workers
WORKER_COUNT.times { queue << :stop }

# Wait for all workers
workers.each(&:join)

puts
puts "✅ All workers finished"
puts

# Alternative: Thread pool with results
puts "=== Thread Pool with Results ==="
puts

class ThreadPool
  def initialize(size)
    @size = size
    @queue = Queue.new
    @workers = []
    start_workers
  end

  def schedule(&block)
    @queue << block
  end

  def shutdown
    @size.times { @queue << :stop }
    @workers.each(&:join)
  end

  private

  def start_workers
    @size.times do |i|
      @workers << Thread.new do
        loop do
          task = @queue.pop
          break if task == :stop
          
          task.call
        end
      end
    end
  end
end

pool = ThreadPool.new(3)

# Schedule tasks
results = []
results_mutex = Mutex.new

5.times do |i|
  pool.schedule do
    sleep 0.1
    result = i * 2
    results_mutex.synchronize do
      results << result
    end
    puts "Task #{i} completed: #{result}"
  end
end

pool.shutdown

puts
puts "Results: #{results.sort}"
puts
puts "✅ Thread pool with results complete"

# Key points:
# - Queue for thread-safe task distribution
# - Fixed number of workers
# - Reusable workers (no thread creation overhead)
# - Manual shutdown needed
# - Good for I/O-bound tasks
