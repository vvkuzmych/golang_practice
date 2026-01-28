package main

import "fmt"

type Mediator interface {
	Notify(sender string, event string)
}

type ChatRoom struct {
	users map[string]*User
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{users: make(map[string]*User)}
}

func (cr *ChatRoom) Register(user *User) {
	cr.users[user.name] = user
	user.mediator = cr
}

func (cr *ChatRoom) Notify(sender string, message string) {
	fmt.Printf("[%s]: %s\n", sender, message)
	for name, user := range cr.users {
		if name != sender {
			user.Receive(sender, message)
		}
	}
}

type User struct {
	name     string
	mediator *ChatRoom
}

func NewUser(name string) *User {
	return &User{name: name}
}

func (u *User) Send(message string) {
	fmt.Printf("%s sends: %s\n", u.name, message)
	u.mediator.Notify(u.name, message)
}

func (u *User) Receive(sender string, message string) {
	fmt.Printf("%s received from %s: %s\n", u.name, sender, message)
}

func main() {
	fmt.Println("=== Mediator Pattern ===\n")

	chatRoom := NewChatRoom()

	alice := NewUser("Alice")
	bob := NewUser("Bob")
	charlie := NewUser("Charlie")

	chatRoom.Register(alice)
	chatRoom.Register(bob)
	chatRoom.Register(charlie)

	alice.Send("Hello everyone!")
	bob.Send("Hi Alice!")
}
