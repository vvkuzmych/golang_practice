package main

import (
	"fmt"
	"time"
)

// Subject - interface для реального об'єкта і proxy
type Subject interface {
	Request(data string) string
}

// RealSubject - реальний об'єкт (дорогий/складний)
type RealSubject struct {
	name string
}

func NewRealSubject(name string) *RealSubject {
	fmt.Printf("Creating RealSubject: %s\n", name)
	// Симуляція дорогої ініціалізації
	time.Sleep(1 * time.Second)
	return &RealSubject{name: name}
}

func (rs *RealSubject) Request(data string) string {
	fmt.Printf("RealSubject handling request: %s\n", data)
	// Симуляція дорогої операції
	time.Sleep(500 * time.Millisecond)
	return fmt.Sprintf("Result from %s: %s", rs.name, data)
}

// Proxy - контролює доступ до RealSubject
type Proxy struct {
	realSubject *RealSubject
	name        string
	cache       map[string]string
}

func NewProxy(name string) *Proxy {
	return &Proxy{
		name:  name,
		cache: make(map[string]string),
	}
}

func (p *Proxy) Request(data string) string {
	// 1. Перевірка кешу
	if cached, ok := p.cache[data]; ok {
		fmt.Printf("Proxy: Returning cached result for: %s\n", data)
		return cached
	}

	// 2. Lazy initialization
	if p.realSubject == nil {
		fmt.Println("Proxy: Creating RealSubject (lazy init)")
		p.realSubject = NewRealSubject(p.name)
	}

	// 3. Додаткова логіка перед запитом
	fmt.Printf("Proxy: Logging request: %s\n", data)

	// 4. Делегування до реального об'єкта
	result := p.realSubject.Request(data)

	// 5. Кешування результату
	p.cache[data] = result

	// 6. Додаткова логіка після запиту
	fmt.Printf("Proxy: Request completed\n")

	return result
}

// ProtectionProxy - контролює права доступу
type ProtectionProxy struct {
	realSubject *RealSubject
	userRole    string
}

func NewProtectionProxy(name string, userRole string) *ProtectionProxy {
	return &ProtectionProxy{
		realSubject: NewRealSubject(name),
		userRole:    userRole,
	}
}

func (pp *ProtectionProxy) Request(data string) string {
	// Перевірка прав доступу
	if pp.userRole != "admin" {
		return "ProtectionProxy: Access denied! Admin role required."
	}

	return pp.realSubject.Request(data)
}

// VirtualProxy - для "дорогих" об'єктів (lazy loading)
type Image interface {
	Display()
}

type RealImage struct {
	filename string
}

func NewRealImage(filename string) *RealImage {
	fmt.Printf("Loading image from disk: %s\n", filename)
	// Симуляція завантаження з диска
	time.Sleep(2 * time.Second)
	return &RealImage{filename: filename}
}

func (ri *RealImage) Display() {
	fmt.Printf("Displaying image: %s\n", ri.filename)
}

type ImageProxy struct {
	filename  string
	realImage *RealImage
}

func NewImageProxy(filename string) *ImageProxy {
	return &ImageProxy{filename: filename}
}

func (ip *ImageProxy) Display() {
	if ip.realImage == nil {
		fmt.Println("ImageProxy: Loading image for the first time...")
		ip.realImage = NewRealImage(ip.filename)
	}
	ip.realImage.Display()
}

func main() {
	fmt.Println("=== 1. Caching Proxy ===")
	proxy := NewProxy("API Service")

	// Перший запит - повільний
	fmt.Println("\nFirst request:")
	result1 := proxy.Request("data1")
	fmt.Printf("Got: %s\n", result1)

	// Другий запит з тими ж даними - швидкий (кеш)
	fmt.Println("\nSecond request (same data):")
	result2 := proxy.Request("data1")
	fmt.Printf("Got: %s\n", result2)

	// Третій запит з новими даними - повільний
	fmt.Println("\nThird request (new data):")
	result3 := proxy.Request("data2")
	fmt.Printf("Got: %s\n", result3)

	fmt.Println("\n=== 2. Protection Proxy ===")
	// User з правами
	adminProxy := NewProtectionProxy("Secure Service", "admin")
	fmt.Println("\nAdmin request:")
	fmt.Println(adminProxy.Request("sensitive data"))

	// User без прав
	userProxy := NewProtectionProxy("Secure Service", "user")
	fmt.Println("\nRegular user request:")
	fmt.Println(userProxy.Request("sensitive data"))

	fmt.Println("\n=== 3. Virtual Proxy (Lazy Loading) ===")
	// Створення proxy - миттєво
	image1 := NewImageProxy("photo1.jpg")
	image2 := NewImageProxy("photo2.jpg")

	fmt.Println("Proxies created (images not loaded yet)")

	// Перше відображення - завантаження
	fmt.Println("\nDisplaying image1 for the first time:")
	image1.Display()

	// Друге відображення - без завантаження
	fmt.Println("\nDisplaying image1 again:")
	image1.Display()

	// Друге зображення
	fmt.Println("\nDisplaying image2:")
	image2.Display()
}
