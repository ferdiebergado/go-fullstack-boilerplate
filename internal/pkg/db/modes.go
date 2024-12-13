package db

type DeleteMode int

const (
	SoftDelete DeleteMode = iota
	HardDelete
)
