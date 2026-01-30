# Week 17 - Ruby Automated Testing

## 🎯 Мета

Автоматизоване тестування: Integration tests, E2E tests, API tests, UI tests.

---

## 📚 Testing Pyramid

```
           /\
          /  \
         / UI \          ← E2E (Capybara)
        /------\
       /  API   \        ← Integration (RSpec)
      /----------\
     /   Unit     \      ← Unit tests (RSpec)
    /--------------\
```

---

## 🛠️ Testing Stack

| Tool | Purpose | Use Case |
|------|---------|----------|
| **RSpec** | Test framework | Unit, integration, API tests |
| **Capybara** | Browser automation | E2E, UI tests |
| **FactoryBot** | Test data | Create test objects |
| **Faker** | Fake data | Generate realistic data |
| **VCR** | HTTP mocking | Record/replay HTTP requests |
| **WebMock** | HTTP stubbing | Stub external APIs |
| **DatabaseCleaner** | DB cleanup | Clean DB between tests |
| **Shoulda Matchers** | RSpec matchers | Simplify Rails tests |

---

## 📖 Theory

### [01: RSpec Basics](./theory/01_rspec_basics.md)

- describe, context, it
- before, after hooks
- let, let!
- Matchers (eq, be, include, raise_error)

### [02: FactoryBot](./theory/02_factory_bot.md)

- Defining factories
- Traits
- Associations
- Sequences

### [03: Capybara E2E](./theory/03_capybara.md)

- Browser automation
- Finding elements
- Interactions (click, fill_in)
- Waiting for elements

### [04: API Testing](./theory/04_api_testing.md)

- Request specs
- JSON parsing
- Response validation
- Authentication

### [05: VCR & WebMock](./theory/05_vcr_webmock.md)

- Recording HTTP interactions
- Stubbing external APIs
- Cassettes

### [06: CI/CD Integration](./theory/06_ci_cd.md)

- GitHub Actions
- GitLab CI
- Parallel testing
- Test coverage

---

## 💻 Practice

### [01: RSpec Basics](./practice/01_rspec_basics/)

- Basic RSpec syntax
- Testing models
- Testing services

### [02: FactoryBot](./practice/02_factory_bot/)

- Creating factories
- Using traits
- Testing associations

### [03: Capybara E2E](./practice/03_capybara_tests/)

- Login flow
- CRUD operations
- Form validation

### [04: API Tests](./practice/04_api_tests/)

- GET/POST/PUT/DELETE
- Authentication
- Error handling

### [05: CI/CD Setup](./practice/05_ci_cd/)

- GitHub Actions config
- Parallel testing
- Coverage reports

---

## 🎯 Test Types

### Unit Tests (Fast)

```ruby
# spec/models/user_spec.rb
RSpec.describe User, type: :model do
  it "validates presence of email" do
    user = User.new(email: nil)
    expect(user).not_to be_valid
  end
end
```

### Integration Tests (Medium)

```ruby
# spec/requests/users_spec.rb
RSpec.describe "Users API", type: :request do
  it "creates a user" do
    post "/api/users", params: { user: { email: "test@example.com" } }
    expect(response).to have_http_status(:created)
  end
end
```

### E2E Tests (Slow)

```ruby
# spec/features/user_login_spec.rb
RSpec.describe "User login", type: :feature do
  it "logs in successfully" do
    visit login_path
    fill_in "Email", with: "user@example.com"
    fill_in "Password", with: "password"
    click_button "Log in"
    expect(page).to have_content("Welcome")
  end
end
```

---

## 🚀 Quick Start

### Installation

```bash
# Gemfile
group :test do
  gem 'rspec-rails'
  gem 'capybara'
  gem 'factory_bot_rails'
  gem 'faker'
  gem 'vcr'
  gem 'webmock'
  gem 'database_cleaner-active_record'
  gem 'shoulda-matchers'
end

# Install
bundle install

# Initialize RSpec
rails generate rspec:install
```

### Run Tests

```bash
# All tests
bundle exec rspec

# Specific file
bundle exec rspec spec/models/user_spec.rb

# Specific line
bundle exec rspec spec/models/user_spec.rb:10

# With format
bundle exec rspec --format documentation

# Parallel
bundle exec parallel_rspec spec/
```

---

## 📊 Best Practices

### 1. Follow AAA Pattern

```ruby
it "creates a user" do
  # Arrange
  user_params = { email: "test@example.com" }
  
  # Act
  post "/api/users", params: user_params
  
  # Assert
  expect(response).to have_http_status(:created)
end
```

### 2. Use FactoryBot

```ruby
# ❌ BAD
user = User.create!(
  email: "test@example.com",
  name: "Test User",
  password: "password123"
)

# ✅ GOOD
user = create(:user)
```

### 3. Use let for Lazy Loading

```ruby
# ✅ GOOD - only created when used
let(:user) { create(:user) }

# ⚠️ let! - created before each test
let!(:admin) { create(:user, :admin) }
```

### 4. Group with describe/context

```ruby
RSpec.describe User do
  describe "#full_name" do
    context "when first and last name present" do
      it "returns full name"
    end
    
    context "when only first name present" do
      it "returns first name"
    end
  end
end
```

### 5. Clean DB Between Tests

```ruby
# spec/rails_helper.rb
RSpec.configure do |config|
  config.use_transactional_fixtures = false
  
  config.before(:suite) do
    DatabaseCleaner.clean_with(:truncation)
  end
  
  config.before(:each) do
    DatabaseCleaner.strategy = :transaction
    DatabaseCleaner.start
  end
  
  config.after(:each) do
    DatabaseCleaner.clean
  end
end
```

---

## 🎯 Testing Checklist

- [ ] **Unit tests** - Models, services, helpers
- [ ] **Integration tests** - API endpoints
- [ ] **E2E tests** - Critical user flows
- [ ] **Test data** - FactoryBot factories
- [ ] **Mocking** - External APIs with VCR
- [ ] **CI/CD** - Automated test runs
- [ ] **Coverage** - 80%+ code coverage
- [ ] **Performance** - Tests run in <5 min

---

## 📚 Quick References

- [RSpec Cheat Sheet](./RSPEC_CHEAT_SHEET.md)
- [Capybara Commands](./CAPYBARA_CHEAT_SHEET.md)
- [FactoryBot Patterns](./FACTORY_BOT_PATTERNS.md)

---

**Week 17: Automated Testing!** 🤖✅
