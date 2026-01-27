# Networking Examples

## TCP Echo Server

### Запуск

Термінал 1 (Server):
```bash
go run tcp_server.go
```

Термінал 2 (Client):
```bash
go run tcp_client.go
```

## Що відбувається?

1. Сервер слухає на порту 9000
2. Клієнт підключається
3. Клієнт відправляє повідомлення
4. Сервер повертає його у верхньому регістрі

## Приклад

```
> hello
Server: HELLO
> world
Server: WORLD
```

## Concepts

- **TCP/IP** - надійний протокол з гарантією доставки
- **net.Listen()** - створює TCP listener
- **net.Dial()** - підключається до TCP сервера
- **Goroutines** - обробка кожного клієнта в окремій goroutine
