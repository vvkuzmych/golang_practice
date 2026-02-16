package main

import (
	"fmt"
	"strconv"
)

// 19. Config - Configuration map pattern

type Config struct {
	values map[string]string
}

func NewConfig() *Config {
	return &Config{values: make(map[string]string)}
}

func (c *Config) Set(key, value string) {
	c.values[key] = value
}

func (c *Config) Get(key string) string {
	return c.values[key]
}

func (c *Config) GetDefault(key, def string) string {
	if v, ok := c.values[key]; ok {
		return v
	}
	return def
}

func (c *Config) GetInt(key string) (int, error) {
	return strconv.Atoi(c.values[key])
}

func (c *Config) GetIntDefault(key string, def int) int {
	if v, err := strconv.Atoi(c.values[key]); err == nil {
		return v
	}
	return def
}

func main() {
	cfg := NewConfig()
	cfg.Set("host", "localhost")
	cfg.Set("port", "8080")
	cfg.Set("debug", "true")

	fmt.Println("Host:", cfg.Get("host"))
	fmt.Println("Missing:", cfg.GetDefault("missing", "default"))
	fmt.Println("Port:", cfg.GetIntDefault("port", 3000))
	fmt.Println("Debug:", cfg.Get("debug"))
}
