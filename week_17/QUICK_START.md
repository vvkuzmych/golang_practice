# Week 17 - Quick Start

## 🚀 Automated Testing Quick Reference

---

## 📦 Installation

```bash
# Gemfile
group :test do
  gem 'rspec-rails'
  gem 'capybara'
  gem 'factory_bot_rails'
  gem 'faker'
  gem 'shoulda-matchers'
  gem 'database_cleaner-active_record'
end

bundle install
rails generate rspec:install
```

---

## 🧪 RSpec Quick Commands

```bash
# Run all tests
bundle exec rspec

# Specific file
bundle exec rspec spec/models/user_spec.rb

# Specific line
bundle exec rspec spec/models/user_spec.rb:10

# With format
bundle exec rspec --format documentation

# Fail fast
bundle exec rspec --fail-fast

# Profile slowest tests
bundle exec rspec --profile 10
```

---

## 📝 RSpec Basic Structure

```ruby
RSpec.describe User do
  describe "#full_name" do
    context "when first and last name present" do
      it "returns full name" do
        user = create(:user, first_name: "John", last_name: "Doe")
        expect(user.full_name).to eq("John Doe")
      end
    end
  end
end
```

---

## 🎭 Common Matchers

```ruby
expect(actual).to eq(expected)
expect(actual).to be_truthy
expect(actual).to be_nil
expect(array).to include(item)
expect { code }.to raise_error(ErrorClass)
expect { code }.to change { User.count }.by(1)
```

---

## 🏭 FactoryBot Quick Usage

```ruby
# Create (saves to DB)
user = create(:user)
users = create_list(:user, 5)

# Build (doesn't save)
user = build(:user)

# With traits
admin = create(:user, :admin)
user_with_posts = create(:user, :with_posts)

# Override attributes
user = create(:user, email: "custom@example.com")
```

---

## 🌐 Capybara Quick Commands

```ruby
# Navigation
visit "/path"
current_path

# Finding
find("#element_id")
all("tr")

# Forms
fill_in "Email", with: "user@example.com"
check "Accept terms"
select "Option", from: "Dropdown"
click_button "Submit"

# Assertions
expect(page).to have_content("Welcome")
expect(page).to have_css(".alert")
expect(page).to have_current_path("/dashboard")
```

---

## 🎯 Test Types

### Unit Test

```ruby
# spec/models/user_spec.rb
RSpec.describe User, type: :model do
  it "validates email" do
    user = build(:user, email: nil)
    expect(user).not_to be_valid
  end
end
```

### API Test

```ruby
# spec/requests/users_spec.rb
RSpec.describe "Users API", type: :request do
  it "creates user" do
    post "/api/users", params: { user: { email: "test@example.com" } }
    expect(response).to have_http_status(:created)
  end
end
```

### E2E Test

```ruby
# spec/features/login_spec.rb
RSpec.describe "Login", type: :feature do
  it "logs in successfully" do
    user = create(:user, password: "password123")
    visit login_path
    fill_in "Email", with: user.email
    fill_in "Password", with: "password123"
    click_button "Log in"
    expect(page).to have_content("Welcome")
  end
end
```

---

## 📚 Files

```bash
cd /Users/vkuzm/GolandProjects/golang_practice/week_17

# Quick references
cat RSPEC_CHEAT_SHEET.md
cat CAPYBARA_CHEAT_SHEET.md
cat FACTORY_BOT_PATTERNS.md
```

---

## ✅ Quick Checklist

- [ ] Встановити gems (RSpec, Capybara, FactoryBot)
- [ ] Створити factories для моделей
- [ ] Написати unit tests для models
- [ ] Написати API tests для endpoints
- [ ] Написати E2E tests для critical flows
- [ ] Налаштувати CI/CD
- [ ] Досягти 80%+ coverage

---

**Week 17: Automated Testing!** 🤖⚡
