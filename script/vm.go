package script

type VM interface {
	Run(script *Script, args []byte) ([]byte, error)
}
