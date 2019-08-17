package types

type Cmd struct {
	Host    string
	Port    int
	Network string
	Genesis struct {
		New   bool
		Coint int
	}
}
