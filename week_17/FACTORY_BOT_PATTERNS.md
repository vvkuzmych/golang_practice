# FactoryBot Patterns

## 🎯 Overview

**FactoryBot** - бібліотека для створення test data в Ruby tests.

---

## 🏗️ Basic Factory

```ruby
# spec/factories/users.rb
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    name { "John Doe" }
    password { "password123" }
  end
end

# Use in tests
user = create(:user)
user = build(:user)                    # Don't save to DB
attributes = attributes_for(:user)      # Get hash of attributes
```

---

## 🎨 Sequences

```ruby
FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "user#{n}@example.com" }
    sequence(:username) { |n| "user#{n}" }
  end
end

# Creates unique emails
create(:user)  # user1@example.com
create(:user)  # user2@example.com
create(:user)  # user3@example.com
```

---

## 🎭 Traits

```ruby
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    name { "John Doe" }
    role { "member" }
    
    trait :admin do
      role { "admin" }
    end
    
    trait :with_posts do
      after(:create) do |user|
        create_list(:post, 3, user: user)
      end
    end
    
    trait :inactive do
      active { false }
    end
  end
end

# Use traits
create(:user, :admin)
create(:user, :admin, :inactive)
create(:user, :with_posts)
```

---

## 🔗 Associations

### Belongs To

```ruby
FactoryBot.define do
  factory :post do
    title { "My Post" }
    body { "Post content" }
    association :user  # Creates associated user
  end
end

# Use
post = create(:post)
post.user  # Associated user created automatically
```

### Has Many

```ruby
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    
    trait :with_posts do
      after(:create) do |user|
        create_list(:post, 3, user: user)
      end
    end
  end
end

# Use
user = create(:user, :with_posts)
user.posts.count  # 3
```

### Nested Associations

```ruby
FactoryBot.define do
  factory :comment do
    body { "Great post!" }
    association :post
    association :user
  end
end

# Creates user, post, and comment
comment = create(:comment)
```

---

## 🎪 Advanced Patterns

### Dynamic Attributes

```ruby
FactoryBot.define do
  factory :user do
    name { Faker::Name.name }
    email { Faker::Internet.email }
    age { rand(18..65) }
    
    # Use Faker for realistic data
    trait :with_address do
      street { Faker::Address.street_address }
      city { Faker::Address.city }
      zip { Faker::Address.zip_code }
    end
  end
end
```

### Dependent Attributes

```ruby
FactoryBot.define do
  factory :user do
    first_name { "John" }
    last_name { "Doe" }
    email { "#{first_name}.#{last_name}@example.com".downcase }
  end
end

create(:user, first_name: "Jane", last_name: "Smith")
# email: jane.smith@example.com
```

### Transient Attributes

```ruby
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    
    transient do
      posts_count { 3 }
    end
    
    after(:create) do |user, evaluator|
      create_list(:post, evaluator.posts_count, user: user)
    end
  end
end

# Use
create(:user, posts_count: 5)
```

---

## 🎯 Factory Inheritance

```ruby
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    name { "John Doe" }
    role { "member" }
    
    factory :admin do
      role { "admin" }
      admin_level { 1 }
    end
    
    factory :super_admin do
      role { "admin" }
      admin_level { 3 }
    end
  end
end

# Use
create(:user)        # member
create(:admin)       # admin with level 1
create(:super_admin) # admin with level 3
```

---

## 📚 Callbacks

```ruby
FactoryBot.define do
  factory :user do
    email { "user@example.com" }
    
    # Before validation
    before(:validation) do |user|
      user.email = user.email.downcase
    end
    
    # After build (not saved)
    after(:build) do |user|
      user.password = "default_password"
    end
    
    # After create (saved)
    after(:create) do |user|
      create(:profile, user: user)
    end
  end
end
```

---

## 🎨 Real-World Examples

### E-commerce System

```ruby
# spec/factories/products.rb
FactoryBot.define do
  factory :product do
    sequence(:name) { |n| "Product #{n}" }
    description { Faker::Lorem.paragraph }
    price { Faker::Commerce.price(range: 10.0..1000.0) }
    stock { rand(0..100) }
    
    trait :out_of_stock do
      stock { 0 }
    end
    
    trait :on_sale do
      sale_price { price * 0.8 }
    end
    
    trait :with_reviews do
      after(:create) do |product|
        create_list(:review, 3, product: product)
      end
    end
  end
  
  factory :order do
    association :user
    status { "pending" }
    total { 0 }
    
    trait :completed do
      status { "completed" }
      completed_at { Time.current }
    end
    
    trait :with_items do
      transient do
        items_count { 2 }
      end
      
      after(:create) do |order, evaluator|
        create_list(:order_item, evaluator.items_count, order: order)
        order.update(total: order.order_items.sum(&:subtotal))
      end
    end
  end
  
  factory :order_item do
    association :order
    association :product
    quantity { rand(1..5) }
    price { product.price }
    
    transient do
      subtotal { quantity * price }
    end
  end
end

# Use
order = create(:order, :with_items, items_count: 5)
product = create(:product, :on_sale, :with_reviews)
```

### Blog System

```ruby
# spec/factories/posts.rb
FactoryBot.define do
  factory :post do
    association :user
    title { Faker::Lorem.sentence }
    body { Faker::Lorem.paragraphs(number: 3).join("\n\n") }
    published { false }
    
    trait :published do
      published { true }
      published_at { Time.current }
    end
    
    trait :with_comments do
      after(:create) do |post|
        create_list(:comment, 5, post: post)
      end
    end
    
    trait :with_tags do
      after(:create) do |post|
        create_list(:tag, 3, posts: [post])
      end
    end
  end
  
  factory :comment do
    association :post
    association :user
    body { Faker::Lorem.paragraph }
    approved { true }
    
    trait :pending do
      approved { false }
    end
  end
  
  factory :tag do
    name { Faker::Lorem.word }
    
    trait :popular do
      transient do
        posts_count { 10 }
      end
      
      after(:create) do |tag, evaluator|
        create_list(:post, evaluator.posts_count, tags: [tag])
      end
    end
  end
end

# Use
post = create(:post, :published, :with_comments, :with_tags)
popular_tag = create(:tag, :popular, posts_count: 20)
```

### Social Network

```ruby
# spec/factories/users.rb
FactoryBot.define do
  factory :user do
    sequence(:email) { |n| "user#{n}@example.com" }
    sequence(:username) { |n| "user#{n}" }
    name { Faker::Name.name }
    bio { Faker::Lorem.sentence }
    
    trait :with_profile do
      after(:create) do |user|
        create(:profile, user: user)
      end
    end
    
    trait :with_followers do
      transient do
        followers_count { 5 }
      end
      
      after(:create) do |user, evaluator|
        evaluator.followers_count.times do
          follower = create(:user)
          create(:follow, follower: follower, following: user)
        end
      end
    end
    
    trait :with_posts do
      transient do
        posts_count { 10 }
      end
      
      after(:create) do |user, evaluator|
        create_list(:post, evaluator.posts_count, user: user)
      end
    end
  end
  
  factory :follow do
    association :follower, factory: :user
    association :following, factory: :user
  end
  
  factory :like do
    association :user
    association :post
  end
end

# Use
user = create(:user, :with_profile, :with_followers, :with_posts, 
              followers_count: 100, posts_count: 50)
```

---

## 📦 Lists

```ruby
# Create multiple records
users = create_list(:user, 5)
users = create_list(:user, 5, :admin)
users = build_list(:user, 5)

# With different attributes
users = create_list(:user, 3) do |user, i|
  user.email = "user#{i}@example.com"
end
```

---

## 🎛️ Configuration

```ruby
# spec/support/factory_bot.rb
RSpec.configure do |config|
  config.include FactoryBot::Syntax::Methods
  
  # Lint factories (check for errors)
  config.before(:suite) do
    FactoryBot.lint
  end
end
```

---

## ✅ Best Practices

### 1. Use Traits

```ruby
# ❌ BAD
create(:user, role: "admin", active: true, verified: true)

# ✅ GOOD
create(:user, :admin, :active, :verified)
```

### 2. Use Sequences for Uniqueness

```ruby
# ❌ BAD
factory :user do
  email { "user@example.com" }  # Will fail on 2nd create
end

# ✅ GOOD
factory :user do
  sequence(:email) { |n| "user#{n}@example.com" }
end
```

### 3. Use Faker for Realistic Data

```ruby
# ❌ BAD
factory :user do
  name { "John Doe" }
end

# ✅ GOOD
factory :user do
  name { Faker::Name.name }
  email { Faker::Internet.email }
end
```

### 4. Keep Factories Simple

```ruby
# ❌ BAD - Too many defaults
factory :user do
  email { "user@example.com" }
  name { "John" }
  age { 25 }
  city { "New York" }
  country { "USA" }
  # ... 20 more attributes
end

# ✅ GOOD - Only required attributes
factory :user do
  sequence(:email) { |n| "user#{n}@example.com" }
  password { "password123" }
  
  trait :complete_profile do
    name { Faker::Name.name }
    age { rand(18..65) }
    # ... other optional attributes
  end
end
```

### 5. Use build Over create When Possible

```ruby
# ❌ SLOWER - Saves to DB
user = create(:user)

# ✅ FASTER - Only builds in memory
user = build(:user)

# Use create when you need DB persistence
user = create(:user)
user.posts.create(...)
```

---

**Week 17: FactoryBot Patterns!** 🏭✅
