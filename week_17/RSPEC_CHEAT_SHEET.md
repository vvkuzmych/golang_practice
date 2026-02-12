# RSpec Cheat Sheet

## 📊 Basic Structure

```ruby
RSpec.describe User do
  describe "#method_name" do
    context "when condition" do
      it "does something" do
        # Test code
      end
    end
  end
end
```

---

## 🎯 Core Methods

### describe / context / it

```ruby
describe "User" do              # Subject
  describe "#full_name" do      # Method
    context "when valid" do     # Condition
      it "returns name" do      # Expectation
        expect(result).to eq("John Doe")
      end
    end
  end
end
```

---

## 🔧 Hooks

```ruby
before(:each)   # Run before each test
before(:all)    # Run once before all tests
after(:each)    # Run after each test
after(:all)     # Run once after all tests

# Example
before(:each) do
  @user = create(:user)
end
```

---

## 💾 let & let!

```ruby
# Lazy evaluation (created when first used)
let(:user) { create(:user) }

# Eager evaluation (created before each test)
let!(:admin) { create(:user, :admin) }

# Use in test
it "works" do
  expect(user.email).to be_present
end
```

---

## ✅ Matchers

### Equality

```ruby
expect(actual).to eq(expected)        # ==
expect(actual).to be(expected)        # equal? (same object)
expect(actual).to eql(expected)       # eql?
expect(actual).not_to eq(expected)    # !=
```

### Truthiness

```ruby
expect(actual).to be_truthy
expect(actual).to be_falsy
expect(actual).to be_nil
expect(actual).to be true
expect(actual).to be false
```

### Comparison

```ruby
expect(actual).to be > expected
expect(actual).to be >= expected
expect(actual).to be < expected
expect(actual).to be <= expected
expect(actual).to be_between(1, 10).inclusive
expect(actual).to be_within(0.1).of(expected)
```

### Collections

```ruby
expect(array).to include(item)
expect(array).to contain_exactly(1, 2, 3)
expect(array).to match_array([3, 2, 1])
expect(hash).to have_key(:key)
expect(hash).to have_value("value")
```

### Types

```ruby
expect(actual).to be_a(String)
expect(actual).to be_an_instance_of(String)
expect(actual).to respond_to(:method_name)
```

### Strings

```ruby
expect(string).to start_with("prefix")
expect(string).to end_with("suffix")
expect(string).to match(/regex/)
```

### Errors

```ruby
expect { code }.to raise_error
expect { code }.to raise_error(ErrorClass)
expect { code }.to raise_error("message")
expect { code }.to raise_error(ErrorClass, "message")
expect { code }.not_to raise_error
```

### Changes

```ruby
expect { code }.to change { object.attribute }
expect { code }.to change { object.attribute }.from(old).to(new)
expect { code }.to change { object.attribute }.by(delta)
expect { code }.to change { User.count }.by(1)
```

---

## 🏗️ Shoulda Matchers (Rails)

### Associations

```ruby
it { should belong_to(:user) }
it { should have_many(:posts) }
it { should have_one(:profile) }
it { should have_and_belong_to_many(:tags) }
```

### Validations

```ruby
it { should validate_presence_of(:email) }
it { should validate_uniqueness_of(:email) }
it { should validate_length_of(:password).is_at_least(8) }
it { should validate_numericality_of(:age) }
it { should validate_inclusion_of(:status).in_array(['active', 'inactive']) }
```

### Database

```ruby
it { should have_db_column(:email).of_type(:string) }
it { should have_db_index(:email) }
```

---

## 🌐 Request Specs (API)

```ruby
RSpec.describe "Users API", type: :request do
  describe "GET /api/users" do
    it "returns users" do
      create_list(:user, 3)
      
      get "/api/users"
      
      expect(response).to have_http_status(:ok)
      expect(JSON.parse(response.body).size).to eq(3)
    end
  end
  
  describe "POST /api/users" do
    it "creates a user" do
      user_params = { user: { email: "test@example.com" } }
      
      post "/api/users", params: user_params
      
      expect(response).to have_http_status(:created)
      expect(User.count).to eq(1)
    end
  end
end
```

---

## 🖥️ Feature Specs (Capybara)

```ruby
RSpec.describe "User login", type: :feature do
  it "logs in successfully" do
    user = create(:user, password: "password123")
    
    visit login_path
    fill_in "Email", with: user.email
    fill_in "Password", with: "password123"
    click_button "Log in"
    
    expect(page).to have_content("Welcome")
    expect(page).to have_current_path(dashboard_path)
  end
end
```

---

## 🎭 Mocking & Stubbing

### Doubles

```ruby
# Create double
user_double = double("User", name: "John", email: "john@example.com")

# Use
expect(user_double.name).to eq("John")
```

### Stubs

```ruby
# Stub method
allow(User).to receive(:find).and_return(user)
allow(user).to receive(:save).and_return(true)

# Stub with arguments
allow(User).to receive(:find).with(1).and_return(user)

# Stub chain
allow(User).to receive_message_chain(:where, :first).and_return(user)
```

### Mocks (Expectations)

```ruby
# Expect method to be called
expect(user).to receive(:save)
expect(user).to receive(:save).with(no_args)
expect(user).to receive(:save).once
expect(user).to receive(:save).twice
expect(user).to receive(:save).exactly(3).times
```

---

## 📝 Shared Examples

```ruby
# Define shared examples
RSpec.shared_examples "a collection" do
  it "is not empty" do
    expect(collection).not_to be_empty
  end
end

# Use shared examples
RSpec.describe Array do
  let(:collection) { [1, 2, 3] }
  it_behaves_like "a collection"
end
```

---

## 🔍 Metadata & Tags

```ruby
# Tag tests
it "works", :focus do
  # Run only this test: rspec --tag focus
end

it "slow test", :slow do
  # Skip: rspec --tag ~slow
end

# Skip test
xit "not implemented yet" do
end

# Pending test
it "will be implemented" do
  pending "waiting for feature"
  expect(true).to be false
end
```

---

## ⚙️ Configuration

```ruby
# spec/spec_helper.rb
RSpec.configure do |config|
  # Use expect syntax (not should)
  config.expect_with :rspec do |expectations|
    expectations.syntax = :expect
  end
  
  # Run tests in random order
  config.order = :random
  
  # Seed for reproducibility
  Kernel.srand config.seed
end
```

---

## 🚀 Quick Commands

```bash
# Run all tests
bundle exec rspec

# Run specific file
bundle exec rspec spec/models/user_spec.rb

# Run specific line
bundle exec rspec spec/models/user_spec.rb:10

# Run with format
bundle exec rspec --format documentation
bundle exec rspec --format html --out results.html

# Run with tags
bundle exec rspec --tag focus
bundle exec rspec --tag ~slow

# Fail fast
bundle exec rspec --fail-fast

# Profile slowest tests
bundle exec rspec --profile 10
```

---

**Week 17: RSpec Cheat Sheet!** 🧪✅
