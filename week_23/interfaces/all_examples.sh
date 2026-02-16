#!/bin/bash
# Generate all interface examples

# 02
cat > 02_empty_interface.go << 'EOF'
package main
import "fmt"
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
func main() {
	describe(42)
	describe("hello")
	describe(true)
	describe([]int{1,2,3})
}
EOF

# 03
cat > 03_type_assertion.go << 'EOF'
package main
import "fmt"
func main() {
	var i interface{} = "hello"
	s := i.(string)
	fmt.Println(s)
	s, ok := i.(string)
	fmt.Println(s, ok)
	f, ok := i.(float64)
	fmt.Println(f, ok)
}
EOF

# 04
cat > 04_type_switch.go << 'EOF'
package main
import "fmt"
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Int: %v\n", v*2)
	case string:
		fmt.Printf("String: %v\n", v+" world")
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}
func main() {
	do(21)
	do("hello")
	do(true)
}
EOF

# 05
cat > 05_multiple_interfaces.go << 'EOF'
package main
import "fmt"
type Reader interface { Read() string }
type Writer interface { Write(string) }
type ReadWriter interface {
	Reader
	Writer
}
type File struct{ data string }
func (f *File) Read() string { return f.data }
func (f *File) Write(s string) { f.data = s }
func main() {
	var rw ReadWriter = &File{}
	rw.Write("Hello")
	fmt.Println(rw.Read())
}
EOF

# 06
cat > 06_stringer.go << 'EOF'
package main
import "fmt"
type Person struct {
	Name string
	Age  int
}
func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}
func main() {
	p := Person{"Alice", 30}
	fmt.Println(p)
}
EOF

# 07
cat > 07_error_interface.go << 'EOF'
package main
import "fmt"
type MyError struct {
	When string
	What string
}
func (e *MyError) Error() string {
	return fmt.Sprintf("%s: %s", e.When, e.What)
}
func run() error {
	return &MyError{"now", "it didn't work"}
}
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
EOF

# 08
cat > 08_polymorphism.go << 'EOF'
package main
import "fmt"
type Shape interface {
	Area() float64
}
type Circle struct{ Radius float64 }
func (c Circle) Area() float64 { return 3.14 * c.Radius * c.Radius }
type Rectangle struct{ Width, Height float64 }
func (r Rectangle) Area() float64 { return r.Width * r.Height }
func printArea(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
}
func main() {
	printArea(Circle{5})
	printArea(Rectangle{3, 4})
}
EOF

# Continue with more examples...
echo "Created 8 interface examples"
