package main

import (
	"fmt"
	"sort"
	"strings"
)

// ============= sort.Interface =============

// type Interface interface {
//     Len() int
//     Less(i, j int) bool
//     Swap(i, j int)
// }

// ============= Person =============

type Person struct {
	Name string
	Age  int
	City string
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d, %s)", p.Name, p.Age, p.City)
}

// ByAge сортує людей за віком
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByName сортує людей за іменем
type ByName []Person

func (n ByName) Len() int           { return len(n) }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

// ============= Product =============

type Product struct {
	Name     string
	Price    float64
	Quantity int
}

func (p Product) String() string {
	return fmt.Sprintf("%s: %.2f UAH (qty: %d)", p.Name, p.Price, p.Quantity)
}

// ByPrice сортує товари за ціною
type ByPrice []Product

func (p ByPrice) Len() int           { return len(p) }
func (p ByPrice) Less(i, j int) bool { return p[i].Price < p[j].Price }
func (p ByPrice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ByQuantity сортує товари за кількістю
type ByQuantity []Product

func (q ByQuantity) Len() int           { return len(q) }
func (q ByQuantity) Less(i, j int) bool { return q[i].Quantity < q[j].Quantity }
func (q ByQuantity) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

// ============= CaseInsensitiveStrings =============

type CaseInsensitiveStrings []string

func (s CaseInsensitiveStrings) Len() int { return len(s) }
func (s CaseInsensitiveStrings) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}
func (s CaseInsensitiveStrings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ============= MultiSort (складне сортування) =============

type MultiSortPerson struct {
	people []Person
	less   []lessFunc
}

type lessFunc func(p1, p2 *Person) bool

func (ms *MultiSortPerson) Sort(people []Person) {
	ms.people = people
	sort.Sort(ms)
}

func (ms *MultiSortPerson) Len() int {
	return len(ms.people)
}

func (ms *MultiSortPerson) Swap(i, j int) {
	ms.people[i], ms.people[j] = ms.people[j], ms.people[i]
}

func (ms *MultiSortPerson) Less(i, j int) bool {
	p, q := &ms.people[i], &ms.people[j]
	for _, less := range ms.less {
		if less(p, q) {
			return true
		}
		if less(q, p) {
			return false
		}
	}
	return false
}

// ============= Reverse Sort =============

type reverse struct {
	sort.Interface
}

func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func Reverse(data sort.Interface) sort.Interface {
	return &reverse{data}
}

// ============= Main =============

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        sort.Interface                    ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ===== Простий приклад з числами =====
	fmt.Println("\n🔹 Сортування чисел (вбудоване)")
	fmt.Println("─────────────────────────────────────────")

	numbers := []int{5, 2, 8, 1, 9, 3}
	fmt.Printf("До:    %v\n", numbers)
	sort.Ints(numbers)
	fmt.Printf("Після: %v\n", numbers)

	// ===== Сортування рядків =====
	fmt.Println("\n🔹 Сортування рядків")
	fmt.Println("─────────────────────────────────────────")

	words := []string{"banana", "apple", "cherry", "date"}
	fmt.Printf("До:    %v\n", words)
	sort.Strings(words)
	fmt.Printf("Після: %v\n", words)

	// ===== Person by Age =====
	fmt.Println("\n🔹 Person - сортування за віком")
	fmt.Println("─────────────────────────────────────────")

	people := []Person{
		{"Іван", 30, "Київ"},
		{"Марія", 25, "Львів"},
		{"Петро", 35, "Одеса"},
		{"Оксана", 28, "Харків"},
	}

	fmt.Println("До сортування:")
	for _, p := range people {
		fmt.Printf("  %s\n", p)
	}

	sort.Sort(ByAge(people))

	fmt.Println("\nПісля сортування за віком:")
	for _, p := range people {
		fmt.Printf("  %s\n", p)
	}

	// ===== Person by Name =====
	fmt.Println("\n🔹 Person - сортування за іменем")
	fmt.Println("─────────────────────────────────────────")

	sort.Sort(ByName(people))

	fmt.Println("Після сортування за іменем:")
	for _, p := range people {
		fmt.Printf("  %s\n", p)
	}

	// ===== Product by Price =====
	fmt.Println("\n🔹 Product - сортування за ціною")
	fmt.Println("─────────────────────────────────────────")

	products := []Product{
		{"Ноутбук", 25000, 5},
		{"Миша", 500, 50},
		{"Клавіатура", 1500, 20},
		{"Монітор", 8000, 10},
	}

	fmt.Println("До сортування:")
	for _, p := range products {
		fmt.Printf("  %s\n", p)
	}

	sort.Sort(ByPrice(products))

	fmt.Println("\nПісля сортування за ціною:")
	for _, p := range products {
		fmt.Printf("  %s\n", p)
	}

	// ===== Product by Quantity =====
	fmt.Println("\n🔹 Product - сортування за кількістю")
	fmt.Println("─────────────────────────────────────────")

	sort.Sort(ByQuantity(products))

	fmt.Println("Після сортування за кількістю:")
	for _, p := range products {
		fmt.Printf("  %s\n", p)
	}

	// ===== Reverse Sort =====
	fmt.Println("\n🔹 Зворотне сортування")
	fmt.Println("─────────────────────────────────────────")

	nums := []int{1, 2, 3, 4, 5}
	fmt.Printf("Оригінал: %v\n", nums)

	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	fmt.Printf("Зворотне: %v\n", nums)

	// Зворотне сортування людей за віком
	sort.Sort(Reverse(ByAge(people)))
	fmt.Println("\nЛюди (від старших до молодших):")
	for _, p := range people {
		fmt.Printf("  %s\n", p)
	}

	// ===== Case-Insensitive Strings =====
	fmt.Println("\n🔹 Case-insensitive сортування")
	fmt.Println("─────────────────────────────────────────")

	mixed := []string{"Zebra", "apple", "Banana", "cherry"}
	fmt.Printf("До:    %v\n", mixed)

	// Звичайне сортування (враховує регістр)
	normalSort := make([]string, len(mixed))
	copy(normalSort, mixed)
	sort.Strings(normalSort)
	fmt.Printf("Звичайне: %v\n", normalSort)

	// Case-insensitive
	sort.Sort(CaseInsensitiveStrings(mixed))
	fmt.Printf("Case-insensitive: %v\n", mixed)

	// ===== sort.Slice (Go 1.8+) =====
	fmt.Println("\n🔹 sort.Slice (зручніший спосіб)")
	fmt.Println("─────────────────────────────────────────")

	people2 := []Person{
		{"Андрій", 28, "Київ"},
		{"Богдан", 32, "Львів"},
		{"Віктор", 25, "Одеса"},
	}

	// Сортування без створення окремого типу
	sort.Slice(people2, func(i, j int) bool {
		return people2[i].Age < people2[j].Age
	})

	fmt.Println("Відсортовано через sort.Slice:")
	for _, p := range people2 {
		fmt.Printf("  %s\n", p)
	}

	// ===== sort.SliceStable (стабільне сортування) =====
	fmt.Println("\n🔹 sort.SliceStable")
	fmt.Println("─────────────────────────────────────────")

	people3 := []Person{
		{"Ігор", 30, "Київ"},
		{"Юлія", 30, "Львів"},
		{"Максим", 30, "Одеса"},
		{"Дарія", 25, "Київ"},
	}

	// Стабільне сортування за віком (порядок однакових зберігається)
	sort.SliceStable(people3, func(i, j int) bool {
		return people3[i].Age < people3[j].Age
	})

	fmt.Println("Стабільне сортування (порядок 30-річних збережено):")
	for _, p := range people3 {
		fmt.Printf("  %s\n", p)
	}

	// ===== Перевірка чи відсортовано =====
	fmt.Println("\n🔹 Перевірка сортування")
	fmt.Println("─────────────────────────────────────────")

	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{5, 2, 8, 1, 9}

	fmt.Printf("%v відсортовано? %t\n", sorted, sort.IntsAreSorted(sorted))
	fmt.Printf("%v відсортовано? %t\n", unsorted, sort.IntsAreSorted(unsorted))

	// ===== Binary Search =====
	fmt.Println("\n🔹 Бінарний пошук")
	fmt.Println("─────────────────────────────────────────")

	sortedNums := []int{1, 3, 5, 7, 9, 11, 13, 15}

	search := 7
	index := sort.SearchInts(sortedNums, search)
	if index < len(sortedNums) && sortedNums[index] == search {
		fmt.Printf("Знайдено %d на позиції %d\n", search, index)
	}

	search = 8
	index = sort.SearchInts(sortedNums, search)
	fmt.Printf("Позиція для вставки %d: %d\n", search, index)

	// ===== Складне сортування =====
	fmt.Println("\n🔹 Складне сортування (кілька критеріїв)")
	fmt.Println("─────────────────────────────────────────")

	people4 := []Person{
		{"Іван", 30, "Київ"},
		{"Марія", 25, "Львів"},
		{"Іван", 25, "Одеса"},
		{"Марія", 30, "Харків"},
	}

	// Спочатку за віком, потім за іменем
	sort.Slice(people4, func(i, j int) bool {
		if people4[i].Age != people4[j].Age {
			return people4[i].Age < people4[j].Age
		}
		return people4[i].Name < people4[j].Name
	})

	fmt.Println("Сортування за віком, потім за іменем:")
	for _, p := range people4 {
		fmt.Printf("  %s\n", p)
	}

	// ===== Summary =====
	fmt.Println("\n\n📝 ВИСНОВКИ")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println("✅ sort.Interface потребує 3 методи:")
	fmt.Println("   • Len() int")
	fmt.Println("   • Less(i, j int) bool")
	fmt.Println("   • Swap(i, j int)")
	fmt.Println()
	fmt.Println("💡 Способи сортування:")
	fmt.Println("   • sort.Sort(data) - власний тип")
	fmt.Println("   • sort.Slice() - lambda функція (зручніше!)")
	fmt.Println("   • sort.SliceStable() - стабільне сортування")
	fmt.Println()
	fmt.Println("🔍 Додаткові можливості:")
	fmt.Println("   • sort.Reverse() - зворотне сортування")
	fmt.Println("   • sort.IsSorted() - перевірка")
	fmt.Println("   • sort.Search() - бінарний пошук")
	fmt.Println()
	fmt.Println("⚡ Рекомендація:")
	fmt.Println("   Використовуйте sort.Slice() замість")
	fmt.Println("   створення окремого типу (простіше!)")
}
