package main

import "fmt"

// State - інтерфейс стану
type State interface {
	InsertCard() string
	EnterPIN(pin string) string
	RequestCash(amount int) string
	EjectCard() string
	GetName() string
}

// ATM - контекст
type ATM struct {
	currentState State
	cardInserted bool
	pinVerified  bool
	balance      int
}

func NewATM(balance int) *ATM {
	atm := &ATM{
		balance: balance,
	}
	atm.currentState = &IdleState{atm: atm}
	return atm
}

func (a *ATM) SetState(state State) {
	fmt.Printf("ATM: Transitioning to %s\n", state.GetName())
	a.currentState = state
}

func (a *ATM) InsertCard() string {
	return a.currentState.InsertCard()
}

func (a *ATM) EnterPIN(pin string) string {
	return a.currentState.EnterPIN(pin)
}

func (a *ATM) RequestCash(amount int) string {
	return a.currentState.RequestCash(amount)
}

func (a *ATM) EjectCard() string {
	return a.currentState.EjectCard()
}

// IdleState - очікування картки
type IdleState struct {
	atm *ATM
}

func (s *IdleState) GetName() string {
	return "IDLE"
}

func (s *IdleState) InsertCard() string {
	s.atm.cardInserted = true
	s.atm.SetState(&CardInsertedState{atm: s.atm})
	return "Card inserted. Please enter PIN."
}

func (s *IdleState) EnterPIN(pin string) string {
	return "Error: No card inserted"
}

func (s *IdleState) RequestCash(amount int) string {
	return "Error: No card inserted"
}

func (s *IdleState) EjectCard() string {
	return "Error: No card inserted"
}

// CardInsertedState - картка вставлена
type CardInsertedState struct {
	atm *ATM
}

func (s *CardInsertedState) GetName() string {
	return "CARD_INSERTED"
}

func (s *CardInsertedState) InsertCard() string {
	return "Error: Card already inserted"
}

func (s *CardInsertedState) EnterPIN(pin string) string {
	if pin == "1234" {
		s.atm.pinVerified = true
		s.atm.SetState(&AuthorizedState{atm: s.atm})
		return "PIN accepted. You may proceed."
	}
	s.atm.SetState(&IdleState{atm: s.atm})
	s.atm.cardInserted = false
	return "Invalid PIN. Card ejected."
}

func (s *CardInsertedState) RequestCash(amount int) string {
	return "Error: Please enter PIN first"
}

func (s *CardInsertedState) EjectCard() string {
	s.atm.cardInserted = false
	s.atm.SetState(&IdleState{atm: s.atm})
	return "Card ejected"
}

// AuthorizedState - авторизовано
type AuthorizedState struct {
	atm *ATM
}

func (s *AuthorizedState) GetName() string {
	return "AUTHORIZED"
}

func (s *AuthorizedState) InsertCard() string {
	return "Error: Card already inserted"
}

func (s *AuthorizedState) EnterPIN(pin string) string {
	return "Error: Already authorized"
}

func (s *AuthorizedState) RequestCash(amount int) string {
	if amount > s.atm.balance {
		return fmt.Sprintf("Error: Insufficient funds. Balance: $%d", s.atm.balance)
	}
	s.atm.balance -= amount
	s.atm.SetState(&DispensingState{atm: s.atm, amount: amount})
	return fmt.Sprintf("Dispensing $%d...", amount)
}

func (s *AuthorizedState) EjectCard() string {
	s.atm.cardInserted = false
	s.atm.pinVerified = false
	s.atm.SetState(&IdleState{atm: s.atm})
	return "Card ejected"
}

// DispensingState - видача грошей
type DispensingState struct {
	atm    *ATM
	amount int
}

func (s *DispensingState) GetName() string {
	return "DISPENSING"
}

func (s *DispensingState) InsertCard() string {
	return "Error: Transaction in progress"
}

func (s *DispensingState) EnterPIN(pin string) string {
	return "Error: Transaction in progress"
}

func (s *DispensingState) RequestCash(amount int) string {
	return "Error: Transaction in progress"
}

func (s *DispensingState) EjectCard() string {
	// Після видачі грошей - повернення до Idle
	s.atm.cardInserted = false
	s.atm.pinVerified = false
	s.atm.SetState(&IdleState{atm: s.atm})
	return fmt.Sprintf("Cash dispensed: $%d. Card ejected. Remaining balance: $%d",
		s.amount, s.atm.balance)
}

// Document Workflow Example

type Document struct {
	state   DocumentState
	content string
	author  string
}

type DocumentState interface {
	Edit(doc *Document, newContent string) string
	Submit(doc *Document) string
	Approve(doc *Document) string
	Reject(doc *Document) string
	Publish(doc *Document) string
	GetName() string
}

// Draft State
type DraftState struct{}

func (s *DraftState) GetName() string { return "DRAFT" }

func (s *DraftState) Edit(doc *Document, newContent string) string {
	doc.content = newContent
	return "Draft updated"
}

func (s *DraftState) Submit(doc *Document) string {
	doc.state = &ReviewState{}
	return "Document submitted for review"
}

func (s *DraftState) Approve(doc *Document) string {
	return "Error: Cannot approve draft"
}

func (s *DraftState) Reject(doc *Document) string {
	return "Error: Cannot reject draft"
}

func (s *DraftState) Publish(doc *Document) string {
	return "Error: Cannot publish draft"
}

// Review State
type ReviewState struct{}

func (s *ReviewState) GetName() string { return "REVIEW" }

func (s *ReviewState) Edit(doc *Document, newContent string) string {
	return "Error: Cannot edit during review"
}

func (s *ReviewState) Submit(doc *Document) string {
	return "Error: Already submitted"
}

func (s *ReviewState) Approve(doc *Document) string {
	doc.state = &ApprovedState{}
	return "Document approved"
}

func (s *ReviewState) Reject(doc *Document) string {
	doc.state = &DraftState{}
	return "Document rejected, back to draft"
}

func (s *ReviewState) Publish(doc *Document) string {
	return "Error: Must be approved first"
}

// Approved State
type ApprovedState struct{}

func (s *ApprovedState) GetName() string { return "APPROVED" }

func (s *ApprovedState) Edit(doc *Document, newContent string) string {
	return "Error: Cannot edit approved document"
}

func (s *ApprovedState) Submit(doc *Document) string {
	return "Error: Already approved"
}

func (s *ApprovedState) Approve(doc *Document) string {
	return "Error: Already approved"
}

func (s *ApprovedState) Reject(doc *Document) string {
	doc.state = &DraftState{}
	return "Approval revoked, back to draft"
}

func (s *ApprovedState) Publish(doc *Document) string {
	doc.state = &PublishedState{}
	return "Document published!"
}

// Published State
type PublishedState struct{}

func (s *PublishedState) GetName() string { return "PUBLISHED" }

func (s *PublishedState) Edit(doc *Document, newContent string) string {
	return "Error: Cannot edit published document"
}

func (s *PublishedState) Submit(doc *Document) string {
	return "Error: Already published"
}

func (s *PublishedState) Approve(doc *Document) string {
	return "Error: Already published"
}

func (s *PublishedState) Reject(doc *Document) string {
	return "Error: Cannot reject published document"
}

func (s *PublishedState) Publish(doc *Document) string {
	return "Error: Already published"
}

func main() {
	fmt.Println("=== 1. ATM State Machine ===\n")

	atm := NewATM(1000)

	// Scenario 1: Успішна транзакція
	fmt.Println("Scenario 1: Successful transaction")
	fmt.Println(atm.InsertCard())
	fmt.Println(atm.EnterPIN("1234"))
	fmt.Println(atm.RequestCash(200))
	fmt.Println(atm.EjectCard())

	fmt.Println("\n" + "=".repeat(50) + "\n")

	// Scenario 2: Невірний PIN
	fmt.Println("Scenario 2: Wrong PIN")
	fmt.Println(atm.InsertCard())
	fmt.Println(atm.EnterPIN("0000"))

	fmt.Println("\n" + "=".repeat(50) + "\n")

	// Scenario 3: Недостатньо коштів
	fmt.Println("Scenario 3: Insufficient funds")
	fmt.Println(atm.InsertCard())
	fmt.Println(atm.EnterPIN("1234"))
	fmt.Println(atm.RequestCash(2000))
	fmt.Println(atm.EjectCard())

	fmt.Println("\n=== 2. Document Workflow ===\n")

	doc := &Document{
		state:   &DraftState{},
		content: "Initial content",
		author:  "John Doe",
	}

	fmt.Printf("Document state: %s\n", doc.state.GetName())
	fmt.Println(doc.state.Edit(doc, "Updated content"))
	fmt.Println(doc.state.Submit(doc))
	fmt.Printf("Document state: %s\n", doc.state.GetName())
	fmt.Println(doc.state.Approve(doc))
	fmt.Printf("Document state: %s\n", doc.state.GetName())
	fmt.Println(doc.state.Publish(doc))
	fmt.Printf("Document state: %s\n", doc.state.GetName())
}
