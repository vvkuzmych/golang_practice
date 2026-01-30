# Capybara Cheat Sheet

## 🎯 Overview

**Capybara** - інструмент для автоматизованого тестування веб-додатків (E2E tests).

---

## 🌐 Navigation

```ruby
visit "/path"                    # Navigate to path
visit root_path                  # Navigate to root
visit user_path(user)            # Navigate with route helper

current_path                     # Get current path
current_url                      # Get current URL

go_back                          # Browser back button
go_forward                       # Browser forward button
```

---

## 🔍 Finding Elements

### By ID/Class/Tag

```ruby
find("#element_id")              # By ID
find(".element_class")           # By class
find("div")                      # By tag

find(:css, "#element_id")        # CSS selector
find(:xpath, "//div[@id='foo']") # XPath
```

### By Text

```ruby
find("a", text: "Click me")      # Exact match
find("a", text: /click/i)        # Regex match
```

### By Attributes

```ruby
find("input[name='email']")
find("button[type='submit']")
find("a[href='/users']")
```

### Multiple Elements

```ruby
all("tr")                        # Find all matching
all("tr").count                  # Count elements
first("tr")                      # First match
all("tr").last                   # Last match
```

---

## 📝 Form Interactions

### Input Fields

```ruby
fill_in "Email", with: "user@example.com"
fill_in "user_email", with: "user@example.com"  # By name
find("#email").set("user@example.com")          # By ID
```

### Checkboxes & Radio Buttons

```ruby
check "Accept terms"
uncheck "Accept terms"
choose "Male"                    # Radio button
```

### Select Dropdowns

```ruby
select "Option", from: "Dropdown"
select "Option", from: "dropdown_id"
unselect "Option", from: "Multi-select"
```

### File Upload

```ruby
attach_file "Avatar", "/path/to/file.jpg"
attach_file "Avatar", Rails.root.join("spec/fixtures/files/avatar.jpg")
```

---

## 🖱️ Click Actions

```ruby
click_link "Sign up"
click_link "link_id"
click_button "Submit"
click_button "submit_button"
click_on "Link or Button"        # Either link or button

# Click specific element
find(".button").click
first(".card").click
```

---

## ✅ Assertions

### Content

```ruby
expect(page).to have_content("Welcome")
expect(page).not_to have_content("Error")
expect(page).to have_text("Hello", exact: true)
expect(page).to have_text(/hello/i)
```

### Elements

```ruby
expect(page).to have_css("div.alert")
expect(page).to have_css("tr", count: 5)
expect(page).to have_selector("table tr", count: 10)
expect(page).to have_xpath("//div[@id='foo']")
```

### Links & Buttons

```ruby
expect(page).to have_link("Sign up")
expect(page).to have_link("Sign up", href: "/signup")
expect(page).to have_button("Submit")
expect(page).to have_button("Submit", disabled: true)
```

### Fields

```ruby
expect(page).to have_field("Email")
expect(page).to have_field("Email", with: "user@example.com")
expect(page).to have_checked_field("Accept terms")
expect(page).to have_unchecked_field("Newsletter")
expect(page).to have_select("Country", selected: "USA")
```

### Current Path

```ruby
expect(page).to have_current_path("/dashboard")
expect(page).to have_current_path(dashboard_path)
expect(page).to have_current_path(/\/users\/\d+/)
```

---

## ⏱️ Waiting

### Default Wait

```ruby
# Capybara автоматично чекає (default: 2 seconds)
expect(page).to have_content("Loaded")  # Waits up to 2s
```

### Custom Wait Time

```ruby
# Change globally
Capybara.default_max_wait_time = 5

# Change per test
using_wait_time(10) do
  expect(page).to have_content("Slow content")
end
```

### Explicit Wait

```ruby
# Wait for element
find("#slow_element", wait: 10)

# Wait for condition
expect(page).to have_css(".loaded", wait: 10)
```

---

## 🪟 Windows & Modals

### Windows/Tabs

```ruby
# Open new window
window = open_new_window
within_window window do
  visit "/new_page"
end

# Switch windows
windows = page.driver.browser.window_handles
page.driver.browser.switch_to.window(windows.last)
```

### Alerts/Confirms

```ruby
# Accept alert
accept_alert do
  click_button "Delete"
end

# Dismiss confirm
dismiss_confirm do
  click_button "Dangerous"
end

# With message check
accept_confirm("Are you sure?") do
  click_button "Delete"
end
```

### Prompts

```ruby
accept_prompt(with: "My input") do
  click_button "Enter name"
end
```

---

## 🎯 Scoping

### Within Block

```ruby
within("#user_form") do
  fill_in "Email", with: "user@example.com"
  click_button "Submit"
end

within("table tbody") do
  expect(page).to have_css("tr", count: 10)
end

within(".card", match: :first) do
  click_link "View"
end
```

### Within Frame

```ruby
within_frame("frame_id") do
  fill_in "Email", with: "user@example.com"
end
```

---

## 🖥️ JavaScript

### Execute JavaScript

```ruby
page.execute_script("alert('Hello')")
page.execute_script("window.scrollTo(0, document.body.scrollHeight)")

result = page.evaluate_script("1 + 1")  # Returns value
```

### Driver

```ruby
# Check if JS driver
Capybara.current_driver == :selenium

# Use JS driver
Capybara.javascript_driver = :selenium_chrome_headless

# In test
it "uses JS", :js do
  # Test with JavaScript support
end
```

---

## 📸 Screenshots & Debugging

```ruby
save_screenshot("screenshot.png")
save_screenshot("tmp/screenshot.png")

save_and_open_screenshot          # Save and open in browser

page.save_page("page.html")       # Save HTML
page.save_and_open_page           # Save and open HTML
```

---

## 🎨 Real-World Examples

### Login Flow

```ruby
RSpec.describe "User login", type: :feature do
  it "logs in successfully" do
    user = create(:user, email: "user@example.com", password: "password123")
    
    visit login_path
    
    fill_in "Email", with: user.email
    fill_in "Password", with: "password123"
    click_button "Log in"
    
    expect(page).to have_content("Welcome, #{user.name}")
    expect(page).to have_current_path(dashboard_path)
  end
end
```

### CRUD Operations

```ruby
RSpec.describe "Post management", type: :feature do
  let(:user) { create(:user) }
  
  before { sign_in user }
  
  it "creates a post" do
    visit new_post_path
    
    fill_in "Title", with: "My Post"
    fill_in "Body", with: "Post content"
    attach_file "Image", Rails.root.join("spec/fixtures/files/image.jpg")
    check "Published"
    select "Technology", from: "Category"
    click_button "Create Post"
    
    expect(page).to have_content("Post created")
    expect(page).to have_content("My Post")
  end
  
  it "edits a post" do
    post = create(:post, user: user, title: "Old Title")
    
    visit edit_post_path(post)
    
    fill_in "Title", with: "New Title"
    click_button "Update Post"
    
    expect(page).to have_content("Post updated")
    expect(page).to have_content("New Title")
    expect(page).not_to have_content("Old Title")
  end
  
  it "deletes a post" do
    post = create(:post, user: user)
    
    visit posts_path
    
    within("#post_#{post.id}") do
      accept_confirm do
        click_link "Delete"
      end
    end
    
    expect(page).to have_content("Post deleted")
    expect(page).not_to have_content(post.title)
  end
end
```

### Form Validation

```ruby
RSpec.describe "Form validation", type: :feature do
  it "shows validation errors" do
    visit new_user_path
    
    click_button "Create User"  # Submit empty form
    
    expect(page).to have_content("Email can't be blank")
    expect(page).to have_content("Password is too short")
    expect(page).to have_css(".field_with_errors")
  end
end
```

### AJAX Interactions

```ruby
RSpec.describe "AJAX interactions", type: :feature, js: true do
  it "loads more items" do
    create_list(:post, 20)
    
    visit posts_path
    
    expect(page).to have_css(".post", count: 10)
    
    click_button "Load More"
    
    expect(page).to have_css(".post", count: 20)
  end
end
```

### Search

```ruby
RSpec.describe "Search", type: :feature do
  it "searches for users" do
    create(:user, name: "John Doe", email: "john@example.com")
    create(:user, name: "Jane Smith", email: "jane@example.com")
    
    visit users_path
    
    fill_in "Search", with: "John"
    click_button "Search"
    
    expect(page).to have_content("John Doe")
    expect(page).not_to have_content("Jane Smith")
  end
end
```

---

## ⚙️ Configuration

```ruby
# spec/rails_helper.rb
RSpec.configure do |config|
  # Use Capybara for feature specs
  config.include Capybara::DSL, type: :feature
  
  # Driver for JS tests
  Capybara.javascript_driver = :selenium_chrome_headless
  
  # Default wait time
  Capybara.default_max_wait_time = 5
  
  # App host (if needed)
  Capybara.app_host = "http://localhost:3000"
end
```

---

## 🚀 Drivers

### Rack::Test (Default)

```ruby
# Fast, no JS support
# spec/rails_helper.rb
Capybara.default_driver = :rack_test
```

### Selenium (Chrome)

```ruby
# spec/rails_helper.rb
Capybara.register_driver :selenium_chrome do |app|
  Capybara::Selenium::Driver.new(app, browser: :chrome)
end

Capybara.javascript_driver = :selenium_chrome
```

### Selenium (Headless Chrome)

```ruby
# spec/rails_helper.rb
Capybara.register_driver :selenium_chrome_headless do |app|
  options = Selenium::WebDriver::Chrome::Options.new
  options.add_argument("--headless")
  options.add_argument("--no-sandbox")
  options.add_argument("--disable-dev-shm-usage")
  
  Capybara::Selenium::Driver.new(app, browser: :chrome, options: options)
end

Capybara.javascript_driver = :selenium_chrome_headless
```

---

## ✅ Best Practices

1. **Use semantic selectors** (text, labels) over CSS/XPath
2. **Wait implicitly** - Capybara waits automatically
3. **Scope interactions** - Use `within` blocks
4. **Test user flows** - Not individual elements
5. **Use data attributes** - `data-testid` for stable selectors
6. **Avoid sleep** - Use Capybara's waiting mechanisms

---

**Week 17: Capybara Cheat Sheet!** 🤖🌐
