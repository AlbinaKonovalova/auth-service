package value

type TokenHash string

func (t TokenHash) String() string {
	return string(t)
}
