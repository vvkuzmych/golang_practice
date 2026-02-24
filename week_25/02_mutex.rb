#!/usr/bin/env ruby

# 02. Mutex in Ruby - Synchronization

puts "=== Ruby Mutex - Synchronization ==="
puts

# Problem: Race condition
puts "1. Race condition (WITHOUT mutex):"
counter = 0
threads = 10.times.map do
  Thread.new do
    1000.times { counter += 1 }
  end
end
threads.each(&:join)
puts "  Expected: 10000"
puts "  Got:      #{counter}"  # Може бути менше через race condition!
puts

# Solution: Mutex
puts "2. Thread-safe (WITH mutex):"
mutex = Mutex.new
counter = 0

threads = 10.times.map do
  Thread.new do
    1000.times do
      mutex.synchronize do
        counter += 1
      end
    end
  end
end

threads.each(&:join)
puts "  Expected: 10000"
puts "  Got:      #{counter}"  # Завжди 10000!
puts

# Mutex methods
puts "3. Mutex methods:"
mutex = Mutex.new

puts "  Locked? #{mutex.locked?}"

mutex.lock
puts "  Locked? #{mutex.locked?}"

mutex.unlock
puts "  After unlock? #{mutex.locked?}"
puts

# Try lock (non-blocking)
puts "4. Try lock (non-blocking):"
mutex = Mutex.new

if mutex.try_lock
  puts "  Lock acquired!"
  mutex.unlock
else
  puts "  Lock failed"
end
puts

# Real-world example: Bank account
puts "5. Real-world: Bank account"
class BankAccount
  def initialize(balance = 0)
    @balance = balance
    @mutex = Mutex.new
  end

  def balance
    @mutex.synchronize { @balance }
  end

  def deposit(amount)
    @mutex.synchronize do
      @balance += amount
    end
  end

  def withdraw(amount)
    @mutex.synchronize do
      return false if @balance < amount
      @balance -= amount
      true
    end
  end
end

account = BankAccount.new(1000)

# Multiple threads accessing account
threads = []
10.times do
  threads << Thread.new { account.deposit(100) }
  threads << Thread.new { account.withdraw(50) }
end

threads.each(&:join)
puts "  Final balance: $#{account.balance}"
puts

puts "✅ Ruby Mutex complete"

# Key points:
# - Mutex.new to create
# - mutex.synchronize { } for safe access
# - mutex.lock / mutex.unlock (manual)
# - mutex.try_lock (non-blocking)
# - Prevents race conditions
