package validecode

type VCodeGenerator interface {
	GetType() string
	Generate() (string, []byte)
}
