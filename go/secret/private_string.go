// Package secret contains a private string that does best effort to prevent the backing string from being accidentally logged.
package secret

type PrivateString string

func (PrivateString) String() string {
	return "private"
}
func (ps PrivateString) AccessPrivateString() string {
	return string(ps)
}
