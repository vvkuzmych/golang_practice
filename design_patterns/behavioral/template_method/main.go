package main

import "fmt"

// DataProcessor - template interface
type DataProcessor interface {
	ReadData() string
	ProcessData(data string) string
	SaveData(data string)
	Execute() // Template Method
}

// BaseProcessor - базова реалізація template method
type BaseProcessor struct {
	processor DataProcessor
}

func (b *BaseProcessor) Execute() {
	data := b.processor.ReadData()
	processed := b.processor.ProcessData(data)
	b.processor.SaveData(processed)
}

// CSVProcessor
type CSVProcessor struct {
	BaseProcessor
	filename string
}

func NewCSVProcessor(filename string) *CSVProcessor {
	p := &CSVProcessor{filename: filename}
	p.BaseProcessor.processor = p
	return p
}

func (p *CSVProcessor) ReadData() string {
	fmt.Printf("CSVProcessor: Reading data from %s\n", p.filename)
	return "csv,data,here"
}

func (p *CSVProcessor) ProcessData(data string) string {
	fmt.Println("CSVProcessor: Processing CSV data")
	return "[PROCESSED_CSV]: " + data
}

func (p *CSVProcessor) SaveData(data string) {
	fmt.Printf("CSVProcessor: Saving to database: %s\n", data)
}

// JSONProcessor
type JSONProcessor struct {
	BaseProcessor
	url string
}

func NewJSONProcessor(url string) *JSONProcessor {
	p := &JSONProcessor{url: url}
	p.BaseProcessor.processor = p
	return p
}

func (p *JSONProcessor) ReadData() string {
	fmt.Printf("JSONProcessor: Fetching data from %s\n", p.url)
	return `{"key": "value"}`
}

func (p *JSONProcessor) ProcessData(data string) string {
	fmt.Println("JSONProcessor: Parsing JSON data")
	return "[PROCESSED_JSON]: " + data
}

func (p *JSONProcessor) SaveData(data string) {
	fmt.Printf("JSONProcessor: Caching result: %s\n", data)
}

func main() {
	fmt.Println("=== Template Method Pattern ===\n")

	csv := NewCSVProcessor("data.csv")
	csv.Execute()

	fmt.Println()

	json := NewJSONProcessor("https://api.example.com/data")
	json.Execute()
}
