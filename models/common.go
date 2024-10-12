package models

type KirbyModel []Kirby

type (
	Kirby struct {
		Info Info
	}
)

type (
	MessageChan struct {
		kch chan KirbyModel
		Ech chan error
	}
)
