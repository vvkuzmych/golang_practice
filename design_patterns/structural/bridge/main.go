package main

import "fmt"

// Implementation - implementation interface
type Device interface {
	IsEnabled() bool
	Enable()
	Disable()
	GetVolume() int
	SetVolume(percent int)
}

// Abstraction - remote control
type Remote struct {
	device Device
}

func NewRemote(device Device) *Remote {
	return &Remote{device: device}
}

func (r *Remote) TogglePower() {
	if r.device.IsEnabled() {
		r.device.Disable()
	} else {
		r.device.Enable()
	}
}

func (r *Remote) VolumeUp() {
	r.device.SetVolume(r.device.GetVolume() + 10)
}

func (r *Remote) VolumeDown() {
	r.device.SetVolume(r.device.GetVolume() - 10)
}

// TV - concrete implementation
type TV struct {
	on     bool
	volume int
}

func (t *TV) IsEnabled() bool { return t.on }
func (t *TV) Enable()         { t.on = true; fmt.Println("TV: ON") }
func (t *TV) Disable()        { t.on = false; fmt.Println("TV: OFF") }
func (t *TV) GetVolume() int  { return t.volume }
func (t *TV) SetVolume(v int) {
	t.volume = v
	fmt.Printf("TV: Volume set to %d\n", v)
}

// Radio - concrete implementation
type Radio struct {
	on     bool
	volume int
}

func (r *Radio) IsEnabled() bool { return r.on }
func (r *Radio) Enable()         { r.on = true; fmt.Println("Radio: ON") }
func (r *Radio) Disable()        { r.on = false; fmt.Println("Radio: OFF") }
func (r *Radio) GetVolume() int  { return r.volume }
func (r *Radio) SetVolume(v int) {
	r.volume = v
	fmt.Printf("Radio: Volume set to %d\n", v)
}

func main() {
	fmt.Println("=== Bridge Pattern ===\n")

	tv := &TV{}
	remote := NewRemote(tv)

	fmt.Println("Testing TV:")
	remote.TogglePower()
	remote.VolumeUp()
	remote.VolumeUp()

	fmt.Println("\nTesting Radio:")
	radio := &Radio{}
	remote.device = radio
	remote.TogglePower()
	remote.VolumeDown()
}
