package execer

import "fmt"

func Example() {
	// Create the execer
	cmd := []string{"echo", "Hello World!"}
	execer, err := Init(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start the command in the background
	err = execer.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get status of the command
	status := execer.Status()
	if status.Error != nil {
		fmt.Println(status.Error)
		return
	}
	fmt.Println(status.Stdout)
}
