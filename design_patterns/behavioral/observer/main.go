package main

import (
	"fmt"
	"strings"
)

// ============= Observer Interface =============

type Observer interface {
	Update(message string)
	GetID() string
}

// ============= Subject (Publisher) =============

type Subject struct {
	observers []Observer
	name      string
}

func NewSubject(name string) *Subject {
	return &Subject{
		observers: []Observer{},
		name:      name,
	}
}

func (s *Subject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
	fmt.Printf("   ➕ %s subscribed to %s\n", observer.GetID(), s.name)
}

func (s *Subject) Detach(observer Observer) {
	for i, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			fmt.Printf("   ➖ %s unsubscribed from %s\n", observer.GetID(), s.name)
			return
		}
	}
}

func (s *Subject) Notify(message string) {
	fmt.Printf("\n📢 %s notifying %d subscribers:\n", s.name, len(s.observers))
	for _, observer := range s.observers {
		observer.Update(message)
	}
}

// ============= Concrete Observers =============

// EmailSubscriber
type EmailSubscriber struct {
	email string
}

func (e *EmailSubscriber) Update(message string) {
	fmt.Printf("   📧 Email to %s: %s\n", e.email, message)
}

func (e *EmailSubscriber) GetID() string {
	return e.email
}

// SMSSubscriber
type SMSSubscriber struct {
	phone string
}

func (s *SMSSubscriber) Update(message string) {
	fmt.Printf("   📱 SMS to %s: %s\n", s.phone, message)
}

func (s *SMSSubscriber) GetID() string {
	return s.phone
}

// MobileAppSubscriber
type MobileAppSubscriber struct {
	userID string
}

func (m *MobileAppSubscriber) Update(message string) {
	fmt.Printf("   📲 Push notification to user %s: %s\n", m.userID, message)
}

func (m *MobileAppSubscriber) GetID() string {
	return m.userID
}

// ============= Example: YouTube Channel =============

type YouTubeChannel struct {
	*Subject
	videoCount int
}

func NewYouTubeChannel(name string) *YouTubeChannel {
	return &YouTubeChannel{
		Subject: NewSubject(name),
	}
}

func (y *YouTubeChannel) UploadVideo(title string) {
	y.videoCount++
	message := fmt.Sprintf("New video: '%s'", title)
	y.Notify(message)
}

// ============= Example: Stock Market =============

type Stock struct {
	*Subject
	symbol string
	price  float64
}

func NewStock(symbol string, price float64) *Stock {
	return &Stock{
		Subject: NewSubject(fmt.Sprintf("Stock %s", symbol)),
		symbol:  symbol,
		price:   price,
	}
}

func (s *Stock) SetPrice(newPrice float64) {
	oldPrice := s.price
	s.price = newPrice
	change := ((newPrice - oldPrice) / oldPrice) * 100
	message := fmt.Sprintf("%s: $%.2f → $%.2f (%.2f%%)",
		s.symbol, oldPrice, newPrice, change)
	s.Notify(message)
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║          Observer Pattern Demo                 ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: YouTube Channel =====
	fmt.Println("\n🔹 Example 1: YouTube Channel Subscriptions")
	fmt.Println(strings.Repeat("─", 50))

	channel := NewYouTubeChannel("Tech Reviews")

	// Subscribers
	sub1 := &EmailSubscriber{email: "user1@example.com"}
	sub2 := &EmailSubscriber{email: "user2@example.com"}
	sub3 := &SMSSubscriber{phone: "+380501234567"}
	sub4 := &MobileAppSubscriber{userID: "user123"}

	fmt.Println("\n📝 Subscribing users:")
	channel.Attach(sub1)
	channel.Attach(sub2)
	channel.Attach(sub3)
	channel.Attach(sub4)

	// Upload video
	channel.UploadVideo("iPhone 15 Pro Review")

	// Unsubscribe
	fmt.Println("\n📝 User unsubscribes:")
	channel.Detach(sub2)

	// Upload another video
	channel.UploadVideo("MacBook Air M3 Unboxing")

	// ===== Example 2: Stock Market =====
	fmt.Println("\n\n🔹 Example 2: Stock Price Notifications")
	fmt.Println(strings.Repeat("─", 50))

	stock := NewStock("AAPL", 150.00)

	// Investors
	investor1 := &EmailSubscriber{email: "investor1@example.com"}
	investor2 := &MobileAppSubscriber{userID: "trader123"}
	investor3 := &SMSSubscriber{phone: "+380501111111"}

	fmt.Println("\n📝 Investors watching stock:")
	stock.Attach(investor1)
	stock.Attach(investor2)
	stock.Attach(investor3)

	// Price changes
	fmt.Println("\n📈 Stock price changes:")
	stock.SetPrice(155.50)
	stock.SetPrice(148.25)

	// ===== Example 3: Multiple Subjects =====
	fmt.Println("\n\n🔹 Example 3: User Subscribed to Multiple Channels")
	fmt.Println(strings.Repeat("─", 50))

	techChannel := NewYouTubeChannel("Tech Channel")
	musicChannel := NewYouTubeChannel("Music Channel")

	user := &EmailSubscriber{email: "multiuser@example.com"}

	fmt.Println("\n📝 User subscribes to multiple channels:")
	techChannel.Attach(user)
	musicChannel.Attach(user)

	fmt.Println("\n📹 Channels upload videos:")
	techChannel.UploadVideo("AI Tutorial")
	musicChannel.UploadVideo("New Song Release")

	// ===== Example 4: Weather Station =====
	fmt.Println("\n\n🔹 Example 4: Weather Station")
	fmt.Println(strings.Repeat("─", 50))

	weatherStation := NewSubject("Weather Station")

	display1 := &MobileAppSubscriber{userID: "phone_display"}
	display2 := &EmailSubscriber{email: "alert@weather.com"}

	fmt.Println("\n📝 Displays register:")
	weatherStation.Attach(display1)
	weatherStation.Attach(display2)

	fmt.Println("\n🌡️  Temperature changes:")
	weatherStation.Notify("Temperature: 25°C, Sunny")
	weatherStation.Notify("Temperature: 18°C, Rainy")

	// ===== Comparison =====
	fmt.Println("\n\n🔹 With vs Without Observer")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Observer (tight coupling):")
	fmt.Println("   youtube.UploadVideo()")
	fmt.Println("   subscriber1.SendEmail()")
	fmt.Println("   subscriber2.SendSMS()")
	fmt.Println("   subscriber3.SendPush()")
	fmt.Println("   → YouTube знає про всіх підписників!")

	fmt.Println("\n✅ With Observer (loose coupling):")
	fmt.Println("   youtube.Notify(message)")
	fmt.Println("   → Підписники самі обробляють повідомлення!")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Слабка зв'язаність (loose coupling)")
	fmt.Println("✅ Динамічні підписки/відписки")
	fmt.Println("✅ Broadcast комунікація (1 → many)")
	fmt.Println("✅ Subject не знає деталей Observer'ів")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Event systems")
	fmt.Println("   - YouTube/Newsletter subscriptions")
	fmt.Println("   - Stock price updates")
	fmt.Println("   - Chat notifications")
	fmt.Println("   - MVC (Model notifies View)")

	fmt.Println("\n🎯 Ключові компоненти:")
	fmt.Println("   Subject (Publisher) - об'єкт що спостерігається")
	fmt.Println("   Observer (Subscriber) - об'єкт що отримує сповіщення")
	fmt.Println("   Update() - метод для отримання сповіщень")

	fmt.Println("\n📚 Go альтернативи:")
	fmt.Println("   - Channels (pub/sub)")
	fmt.Println("   - context.Context cancellation")
}
