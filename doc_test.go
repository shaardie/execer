package execer

func ExampleNewExecer() {
	cmd := []string{"echo", "Hello World!"}
	execer := Init(cmd)
}
