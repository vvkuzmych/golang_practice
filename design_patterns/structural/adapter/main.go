package main

import (
	"fmt"
	"strings"
)

// ============= Example 1: Media Player Adapter =============

// MediaPlayer - цільовий інтерфейс
type MediaPlayer interface {
	Play(filename string) string
}

// MP3Player - працює з MP3
type MP3Player struct{}

func (m *MP3Player) Play(filename string) string {
	return fmt.Sprintf("🎵 Playing MP3: %s", filename)
}

// VLCPlayer - legacy player (несумісний інтерфейс)
type VLCPlayer struct{}

func (v *VLCPlayer) PlayVLC(filename string) string {
	return fmt.Sprintf("🎬 VLC playing: %s", filename)
}

// MP4Player - інша стороння бібліотека
type MP4Player struct{}

func (m *MP4Player) PlayMP4File(filename string) string {
	return fmt.Sprintf("📹 MP4 player: %s", filename)
}

// VLCAdapter - адаптер для VLC
type VLCAdapter struct {
	vlc *VLCPlayer
}

func (a *VLCAdapter) Play(filename string) string {
	return a.vlc.PlayVLC(filename)
}

// MP4Adapter - адаптер для MP4
type MP4Adapter struct {
	mp4 *MP4Player
}

func (a *MP4Adapter) Play(filename string) string {
	return a.mp4.PlayMP4File(filename)
}

// ============= Example 2: Payment Gateway Adapter =============

// PaymentProcessor - наш інтерфейс
type PaymentProcessor interface {
	ProcessPayment(amount float64) string
	GetFee() float64
}

// StripeAPI - сторонній API (інший інтерфейс)
type StripeAPI struct{}

func (s *StripeAPI) Charge(cents int) string {
	return fmt.Sprintf("💳 Stripe charged %d cents", cents)
}

func (s *StripeAPI) GetStripeFee() int {
	return 30 // cents
}

// PayPalAPI - інший API
type PayPalAPI struct{}

func (p *PayPalAPI) SendPayment(dollars float64) string {
	return fmt.Sprintf("💰 PayPal sent $%.2f", dollars)
}

func (p *PayPalAPI) PayPalCommission() float64 {
	return 0.029 // 2.9%
}

// StripeAdapter
type StripeAdapter struct {
	stripe *StripeAPI
}

func (a *StripeAdapter) ProcessPayment(amount float64) string {
	cents := int(amount * 100)
	return a.stripe.Charge(cents)
}

func (a *StripeAdapter) GetFee() float64 {
	return float64(a.stripe.GetStripeFee()) / 100.0
}

// PayPalAdapter
type PayPalAdapter struct {
	paypal *PayPalAPI
}

func (a *PayPalAdapter) ProcessPayment(amount float64) string {
	return a.paypal.SendPayment(amount)
}

func (a *PayPalAdapter) GetFee() float64 {
	return a.paypal.PayPalCommission()
}

// ============= Example 3: Temperature Converter Adapter =============

// TemperatureReader - наш інтерфейс (Celsius)
type TemperatureReader interface {
	ReadCelsius() float64
	GetLocation() string
}

// FahrenheitSensor - legacy sensor
type FahrenheitSensor struct {
	location string
	temp     float64
}

func (f *FahrenheitSensor) ReadFahrenheit() float64 {
	return f.temp
}

func (f *FahrenheitSensor) GetSensorLocation() string {
	return f.location
}

// FahrenheitAdapter
type FahrenheitAdapter struct {
	sensor *FahrenheitSensor
}

func (a *FahrenheitAdapter) ReadCelsius() float64 {
	f := a.sensor.ReadFahrenheit()
	return (f - 32) * 5 / 9
}

func (a *FahrenheitAdapter) GetLocation() string {
	return a.sensor.GetSensorLocation()
}

// ============= Main =============

func main() {
	fmt.Println("╔════════════════════════════════════════════════╗")
	fmt.Println("║          Adapter Pattern Demo                  ║")
	fmt.Println("╚════════════════════════════════════════════════╝")

	// ===== Example 1: Media Player =====
	fmt.Println("\n🔹 Example 1: Media Player Adapters")
	fmt.Println(strings.Repeat("─", 50))

	// Native player
	var player MediaPlayer = &MP3Player{}
	fmt.Println(player.Play("song.mp3"))

	// VLC через адаптер
	player = &VLCAdapter{vlc: &VLCPlayer{}}
	fmt.Println(player.Play("movie.avi"))

	// MP4 через адаптер
	player = &MP4Adapter{mp4: &MP4Player{}}
	fmt.Println(player.Play("video.mp4"))

	fmt.Println("\n✅ Всі плеєри працюють через єдиний інтерфейс!")

	// ===== Example 2: Payment Gateways =====
	fmt.Println("\n\n🔹 Example 2: Payment Gateway Adapters")
	fmt.Println(strings.Repeat("─", 50))

	amount := 100.00

	// Stripe через адаптер
	var processor PaymentProcessor = &StripeAdapter{stripe: &StripeAPI{}}
	fmt.Printf("\n💳 Stripe:\n")
	fmt.Printf("   %s\n", processor.ProcessPayment(amount))
	fmt.Printf("   Fee: $%.2f\n", processor.GetFee())

	// PayPal через адаптер
	processor = &PayPalAdapter{paypal: &PayPalAPI{}}
	fmt.Printf("\n💰 PayPal:\n")
	fmt.Printf("   %s\n", processor.ProcessPayment(amount))
	fmt.Printf("   Fee: %.1f%%\n", processor.GetFee()*100)

	fmt.Println("\n✅ Різні API, але єдиний інтерфейс!")

	// ===== Example 3: Temperature Sensors =====
	fmt.Println("\n\n🔹 Example 3: Temperature Sensor Adapter")
	fmt.Println(strings.Repeat("─", 50))

	// Legacy Fahrenheit sensor
	fahrenheitSensor := &FahrenheitSensor{
		location: "New York",
		temp:     77.0, // °F
	}

	// Адаптер для конвертації в Celsius
	var tempReader TemperatureReader = &FahrenheitAdapter{
		sensor: fahrenheitSensor,
	}

	fmt.Printf("\n🌡️  Location: %s\n", tempReader.GetLocation())
	fmt.Printf("   Temperature: %.1f°F → %.1f°C\n",
		fahrenheitSensor.ReadFahrenheit(),
		tempReader.ReadCelsius())

	// ===== Example 4: Multiple Adapters =====
	fmt.Println("\n\n🔹 Example 4: Using Multiple Players")
	fmt.Println(strings.Repeat("─", 50))

	playlist := []struct {
		file   string
		player MediaPlayer
	}{
		{"song1.mp3", &MP3Player{}},
		{"movie.avi", &VLCAdapter{vlc: &VLCPlayer{}}},
		{"video.mp4", &MP4Adapter{mp4: &MP4Player{}}},
		{"song2.mp3", &MP3Player{}},
	}

	fmt.Println("\n🎵 Playing playlist:")
	for i, item := range playlist {
		fmt.Printf("   %d. %s\n", i+1, item.player.Play(item.file))
	}

	// ===== Comparison =====
	fmt.Println("\n\n🔹 Without vs With Adapter")
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println("\n❌ Without Adapter:")
	fmt.Println("   if type == MP3 { mp3.Play() }")
	fmt.Println("   if type == VLC { vlc.PlayVLC() }")
	fmt.Println("   if type == MP4 { mp4.PlayMP4File() }")
	fmt.Println("   → Різні методи, складний код!")

	fmt.Println("\n✅ With Adapter:")
	fmt.Println("   player.Play() // для всіх!")
	fmt.Println("   → Єдиний інтерфейс, простий код!")

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("✅ Працює з несумісними інтерфейсами")
	fmt.Println("✅ Не змінює існуючий код")
	fmt.Println("✅ Інтеграція legacy систем")
	fmt.Println("✅ Обгортка сторонніх бібліотек")

	fmt.Println("\n💡 ВИКОРИСТАННЯ:")
	fmt.Println("   - Legacy code integration")
	fmt.Println("   - Third-party API wrappers")
	fmt.Println("   - Data format converters")
	fmt.Println("   - Incompatible interface bridging")

	fmt.Println("\n📚 Реальні приклади:")
	fmt.Println("   - USB-C → USB-A адаптер")
	fmt.Println("   - 220V → 110V трансформатор")
	fmt.Println("   - API v1 → API v2 wrapper")
}
