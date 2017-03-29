package execer

func Example() {
	cmd := []string{"echo", "Hello World!"}
	_, err := Init(cmd)
	if err != nil {
		// Something went wrong during creation
	}
}
