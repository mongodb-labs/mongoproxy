package server

var Registry = make(map[string]Module)

func Publish(m Module) {
	Registry[m.Name()] = m
}
