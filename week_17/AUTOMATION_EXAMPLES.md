# Automation Testing Examples for Beginners

## 🎯 Що таке Automation Testing?

**Automation Testing** - автоматизоване тестування веб-додатків через браузер:
- Відкриває браузер
- Виконує дії користувача (клік, введення тексту)
- Перевіряє результати
- Закриває браузер

**Інструменти:** Capybara + Selenium

---

## 📦 Setup (Крок за кроком)

### 1. Gemfile

```ruby
group :test do
  gem 'capybara'
  gem 'selenium-webdriver'
  gem 'rspec'
end
```

```bash
bundle install
```

### 2. spec/spec_helper.rb

```ruby
require 'capybara/rspec'
require 'selenium-webdriver'

# Налаштування Capybara
Capybara.register_driver :selenium_chrome do |app|
  options = Selenium::WebDriver::Chrome::Options.new
  options.add_argument('--headless')  # Без вікна браузера
  options.add_argument('--no-sandbox')
  options.add_argument('--disable-dev-shm-usage')
  
  Capybara::Selenium::Driver.new(app, browser: :chrome, options: options)
end

Capybara.default_driver = :selenium_chrome
Capybara.javascript_driver = :selenium_chrome

# Час очікування
Capybara.default_max_wait_time = 5

RSpec.configure do |config|
  config.include Capybara::DSL
end
```

### 3. Запуск

```bash
bundle exec rspec spec/features/
```

---

## 🌐 Example 1: Відкрити сайт і перевірити title

```ruby
# spec/features/google_search_spec.rb
require 'spec_helper'

describe "Google Search", type: :feature do
  it "відкриває Google і перевіряє title" do
    # 1. Відкрити сайт
    visit "https://www.google.com"
    
    # 2. Перевірити title
    expect(page).to have_title("Google")
    
    # 3. Перевірити, що є поле пошуку
    expect(page).to have_css("textarea[name='q']")
    
    puts "✅ Google відкрито успішно!"
  end
end
```

**Запуск:**
```bash
bundle exec rspec spec/features/google_search_spec.rb
```

---

## 🔍 Example 2: Пошук на Google

```ruby
# spec/features/google_search_spec.rb
require 'spec_helper'

describe "Google Search", type: :feature do
  it "виконує пошук на Google" do
    # 1. Відкрити Google
    visit "https://www.google.com"
    
    # 2. Знайти поле пошуку
    search_box = find("textarea[name='q']")
    
    # 3. Ввести текст
    search_box.set("Capybara automation testing")
    
    # 4. Натиснути Enter
    search_box.send_keys(:return)
    
    # 5. Дочекатися результатів
    expect(page).to have_content("results")
    
    # 6. Перевірити, що є результати
    expect(page).to have_css("#search")
    
    puts "✅ Пошук виконано успішно!"
  end
end
```

---

## 📝 Example 3: Заповнити форму (Registration)

```ruby
# spec/features/registration_spec.rb
require 'spec_helper'

describe "User Registration", type: :feature do
  it "заповнює форму реєстрації" do
    # 1. Відкрити сторінку реєстрації
    visit "https://example.com/signup"  # Замініть на ваш URL
    
    # 2. Заповнити форму
    fill_in "Name", with: "John Doe"
    fill_in "Email", with: "john@example.com"
    fill_in "Password", with: "SecurePass123"
    fill_in "Password Confirmation", with: "SecurePass123"
    
    # 3. Вибрати checkbox
    check "I accept terms and conditions"
    
    # 4. Вибрати з dropdown
    select "United States", from: "Country"
    
    # 5. Вибрати radio button
    choose "Male"
    
    # 6. Натиснути кнопку
    click_button "Sign Up"
    
    # 7. Перевірити успіх
    expect(page).to have_content("Welcome")
    expect(page).to have_content("Registration successful")
    
    puts "✅ Реєстрація успішна!"
  end
end
```

---

## 🔐 Example 4: Login Flow

```ruby
# spec/features/login_spec.rb
require 'spec_helper'

describe "User Login", type: :feature do
  it "логінить користувача" do
    # 1. Відкрити login сторінку
    visit "https://example.com/login"
    
    # 2. Знайти поля за ID
    find("#email").set("user@example.com")
    find("#password").set("password123")
    
    # 3. Або заповнити за label
    fill_in "Email", with: "user@example.com"
    fill_in "Password", with: "password123"
    
    # 4. Натиснути Login
    click_button "Log In"
    
    # 5. Дочекатися редіректу
    expect(page).to have_current_path("/dashboard")
    
    # 6. Перевірити, що залогінено
    expect(page).to have_content("Welcome back")
    
    # 7. Перевірити, що є кнопка Logout
    expect(page).to have_link("Logout")
    
    puts "✅ Логін успішний!"
  end
  
  it "показує помилку при невірному паролі" do
    visit "https://example.com/login"
    
    fill_in "Email", with: "user@example.com"
    fill_in "Password", with: "wrong_password"
    click_button "Log In"
    
    # Перевірити помилку
    expect(page).to have_content("Invalid email or password")
    expect(page).to have_css(".alert-danger")
    
    puts "✅ Помилка показана коректно!"
  end
end
```

---

## 🛒 Example 5: E-commerce Flow (Add to Cart)

```ruby
# spec/features/shopping_spec.rb
require 'spec_helper'

describe "Shopping Cart", type: :feature do
  it "додає товар в кошик" do
    # 1. Відкрити головну сторінку
    visit "https://example-shop.com"
    
    # 2. Знайти товар
    within(".product-list") do
      click_link "View Product", match: :first
    end
    
    # 3. Вибрати розмір
    select "Large", from: "Size"
    
    # 4. Вибрати колір
    find(".color-option[data-color='blue']").click
    
    # 5. Додати в кошик
    click_button "Add to Cart"
    
    # 6. Перевірити notification
    expect(page).to have_content("Item added to cart")
    
    # 7. Перейти в кошик
    click_link "Cart"
    
    # 8. Перевірити товар в кошику
    expect(page).to have_content("Large")
    expect(page).to have_content("Blue")
    
    # 9. Перевірити кількість
    expect(page).to have_css(".cart-items", count: 1)
    
    puts "✅ Товар додано в кошик!"
  end
end
```

---

## 📸 Example 6: Завантажити файл (Upload)

```ruby
# spec/features/file_upload_spec.rb
require 'spec_helper'

describe "File Upload", type: :feature do
  it "завантажує файл" do
    # 1. Відкрити сторінку
    visit "https://example.com/upload"
    
    # 2. Вибрати файл
    attach_file("Avatar", Rails.root.join("spec/fixtures/files/avatar.jpg"))
    
    # Або якщо немає Rails:
    # attach_file("Avatar", File.expand_path("../fixtures/avatar.jpg", __FILE__))
    
    # 3. Додати опис
    fill_in "Description", with: "My profile picture"
    
    # 4. Завантажити
    click_button "Upload"
    
    # 5. Перевірити успіх
    expect(page).to have_content("File uploaded successfully")
    
    # 6. Перевірити, що зображення відображається
    expect(page).to have_css("img[src*='avatar']")
    
    puts "✅ Файл завантажено!"
  end
end
```

---

## ⚠️ Example 7: Обробка Alert/Confirm

```ruby
# spec/features/alerts_spec.rb
require 'spec_helper'

describe "Alerts and Confirms", type: :feature do
  it "обробляє JavaScript alert" do
    visit "https://example.com/delete"
    
    # Прийняти alert
    accept_alert do
      click_button "Delete Account"
    end
    
    expect(page).to have_content("Account deleted")
  end
  
  it "обробляє confirm dialog" do
    visit "https://example.com/settings"
    
    # Прийняти confirm
    accept_confirm("Are you sure?") do
      click_button "Reset Settings"
    end
    
    expect(page).to have_content("Settings reset")
  end
  
  it "відхиляє confirm dialog" do
    visit "https://example.com/settings"
    
    # Відхилити confirm
    dismiss_confirm do
      click_button "Reset Settings"
    end
    
    expect(page).not_to have_content("Settings reset")
  end
end
```

---

## ⏱️ Example 8: Очікування елементів (Wait for AJAX)

```ruby
# spec/features/ajax_spec.rb
require 'spec_helper'

describe "AJAX Loading", type: :feature do
  it "чекає завантаження даних" do
    visit "https://example.com/dashboard"
    
    # Натиснути кнопку, яка завантажує дані
    click_button "Load More"
    
    # Дочекатися появи нових елементів
    expect(page).to have_css(".item", count: 20, wait: 10)
    
    # Дочекатися зникнення loader
    expect(page).not_to have_css(".loading-spinner")
    
    # Дочекатися конкретного тексту
    expect(page).to have_content("Loaded", wait: 5)
  end
  
  it "перевіряє динамічний контент" do
    visit "https://example.com/search"
    
    # Ввести в поле пошуку
    fill_in "Search", with: "Ruby"
    
    # Дочекатися автокомпліту
    expect(page).to have_css(".autocomplete-results", wait: 5)
    
    # Вибрати перший результат
    within(".autocomplete-results") do
      first("li").click
    end
    
    expect(page).to have_content("Search results for: Ruby")
  end
end
```

---

## 📱 Example 9: Multiple Windows/Tabs

```ruby
# spec/features/windows_spec.rb
require 'spec_helper'

describe "Multiple Windows", type: :feature do
  it "відкриває нове вікно і працює з ним" do
    visit "https://example.com"
    
    # Відкрити нове вікно (клік на посилання з target="_blank")
    click_link "Open in New Tab"
    
    # Перемкнутися на нове вікно
    new_window = window_opened_by { click_link "Open in New Tab" }
    
    within_window new_window do
      expect(page).to have_content("New Window Content")
      
      # Виконати дії в новому вікні
      fill_in "Search", with: "Test"
      click_button "Submit"
    end
    
    # Повернутися в основне вікно
    # (автоматично після within_window block)
    expect(page).to have_content("Original Window Content")
  end
end
```

---

## 🎯 Example 10: Реальний повний сценарій

```ruby
# spec/features/complete_user_journey_spec.rb
require 'spec_helper'

describe "Complete User Journey", type: :feature do
  it "проходить повний шлях користувача" do
    # 1. Відкрити головну сторінку
    visit "https://example.com"
    expect(page).to have_content("Welcome")
    puts "✅ Крок 1: Головна сторінка"
    
    # 2. Перейти на реєстрацію
    click_link "Sign Up"
    expect(page).to have_current_path("/signup")
    puts "✅ Крок 2: Сторінка реєстрації"
    
    # 3. Заповнити форму реєстрації
    fill_in "Name", with: "Test User"
    fill_in "Email", with: "test#{Time.now.to_i}@example.com"
    fill_in "Password", with: "SecurePass123"
    fill_in "Password Confirmation", with: "SecurePass123"
    check "I accept terms"
    click_button "Create Account"
    puts "✅ Крок 3: Форма заповнена"
    
    # 4. Перевірити успішну реєстрацію
    expect(page).to have_content("Registration successful")
    expect(page).to have_current_path("/dashboard")
    puts "✅ Крок 4: Реєстрація успішна"
    
    # 5. Перейти в профіль
    click_link "Profile"
    expect(page).to have_content("Test User")
    puts "✅ Крок 5: Профіль відкрито"
    
    # 6. Редагувати профіль
    click_link "Edit Profile"
    fill_in "Bio", with: "This is my bio"
    attach_file("Avatar", File.expand_path("../fixtures/avatar.jpg", __FILE__))
    click_button "Save Changes"
    puts "✅ Крок 6: Профіль відредаговано"
    
    # 7. Перевірити зміни
    expect(page).to have_content("Profile updated")
    expect(page).to have_content("This is my bio")
    expect(page).to have_css("img[alt='Avatar']")
    puts "✅ Крок 7: Зміни збережено"
    
    # 8. Logout
    click_link "Logout"
    expect(page).to have_content("Logged out")
    expect(page).to have_current_path("/")
    puts "✅ Крок 8: Logout успішний"
    
    puts "\n🎉 ПОВНИЙ СЦЕНАРІЙ ПРОЙДЕНО УСПІШНО!"
  end
end
```

---

## 📊 Example 11: Скріншоти при помилках

```ruby
# spec/spec_helper.rb
RSpec.configure do |config|
  # Зробити скріншот при падінні тесту
  config.after(:each, type: :feature) do |example|
    if example.exception
      screenshot_name = "screenshot_#{Time.now.to_i}.png"
      save_screenshot(screenshot_name)
      puts "💾 Screenshot saved: #{screenshot_name}"
    end
  end
end
```

```ruby
# spec/features/example_spec.rb
describe "Test with screenshots", type: :feature do
  it "робить скріншот при помилці" do
    visit "https://example.com"
    
    # Якщо тест впаде, автоматично збережеться скріншот
    expect(page).to have_content("Non-existent content")
  end
  
  it "робить скріншот вручну" do
    visit "https://example.com/dashboard"
    
    # Зробити скріншот вручну
    save_screenshot("dashboard.png")
    
    # Або зробити і відкрити
    save_and_open_screenshot
  end
end
```

---

## 🎨 Example 12: Selenium з видимим браузером (для дебагу)

```ruby
# spec/spec_helper.rb

# Для дебагу - показувати браузер
Capybara.register_driver :selenium_chrome_visible do |app|
  options = Selenium::WebDriver::Chrome::Options.new
  # НЕ додаємо --headless, щоб бачити браузер
  
  Capybara::Selenium::Driver.new(app, browser: :chrome, options: options)
end

# Використовувати для конкретного тесту:
describe "Debug test", type: :feature, driver: :selenium_chrome_visible do
  it "показує браузер" do
    visit "https://google.com"
    
    # Затримка, щоб побачити браузер
    sleep 5
    
    fill_in "q", with: "Test"
    sleep 2
    
    find("input[name='q']").send_keys(:return)
    sleep 5
  end
end
```

---

## ⚡ Корисні команди

```bash
# Запустити всі тести
bundle exec rspec spec/features/

# Запустити конкретний файл
bundle exec rspec spec/features/login_spec.rb

# Запустити конкретний тест (по рядку)
bundle exec rspec spec/features/login_spec.rb:10

# З деталями
bundle exec rspec spec/features/ --format documentation

# Fail fast (зупинитися на першій помилці)
bundle exec rspec spec/features/ --fail-fast
```

---

## 🎯 Tips для початківців

### 1. Використовуй inspect для дебагу

```ruby
it "дебагу елемент" do
  visit "https://example.com"
  
  # Подивитися HTML сторінки
  puts page.html
  
  # Знайти всі кнопки
  puts all("button").map(&:text)
  
  # Знайти всі посилання
  puts all("a").map { |link| link[:href] }
end
```

### 2. Використовуй save_and_open_page

```ruby
it "дебагу сторінку" do
  visit "https://example.com"
  fill_in "Email", with: "test@example.com"
  
  # Зберегти HTML і відкрити в браузері
  save_and_open_page
end
```

### 3. Використовуй data-test-id

```html
<!-- HTML -->
<button data-test-id="submit-button">Submit</button>
```

```ruby
# Test
find("[data-test-id='submit-button']").click
```

### 4. Чекай елементи (не використовуй sleep)

```ruby
# ❌ BAD
sleep 5

# ✅ GOOD
expect(page).to have_content("Loaded", wait: 10)
```

---

**Week 17: Automation Testing Examples!** 🤖✅
