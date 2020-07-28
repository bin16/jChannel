package main

type channel struct {
	ID       int64
	Username string
}

func (ch channel) Recipient() string {
	return ch.Username
}
