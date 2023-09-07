package secret

type PrivateString string

func (PrivateString) String() string {
	return "private"
}
func (ps PrivateString) AccessPrivateString() string {
	return string(ps)
}
