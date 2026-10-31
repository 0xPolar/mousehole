package crypt

type ISigner interface {
	Sign()
	Save()
	Load()
	Verify()
}
