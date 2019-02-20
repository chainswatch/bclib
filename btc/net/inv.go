package net

// Hash returns inventory hash
func (i *Inv) Hash() [32]byte {
	return i.hash
}

// Object returns inventory object
func (i *Inv) Object() string {
	return i.object
}

// Payload returns inventory payload
func (i *Inv) Payload() []byte {
	return i.payload
}
