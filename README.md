# GO Terminator - Graceful Process Termination Utility


[![Go Report Card](https://goreportcard.com/badge/github.com/RohanPoojary/go-terminator)](https://goreportcard.com/report/github.com/RohanPoojary/go-terminator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/RohanPoojary/go-terminator/blob/main/LICENSE)


The `go-terminator` package provides a utility in Go for gracefully terminating processes by closing registered resources when a termination signal is received.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
  - [Creating a Terminator](#creating-a-terminator)
  - [Adding Resources](#adding-resources)
  - [Setting Callback](#setting-callback)
  - [Waiting for Termination](#waiting-for-termination)
- [Termination Result Structure](#terminationresult-structure)
- [Contributing](#contributing)
- [License](#license)

## Introduction

When building long-running applications or services, it's important to ensure that resources are properly cleaned up when the process is terminated. The `go-terminator` package helps manage the graceful termination of your application by providing a simple mechanism to register resources that need to be closed when the process receives termination signals.

## Installation

To use the `terminator` package, import it in your Go code:

```shell
go get "github.com/RohanPoojary/go-terminator"
```

## Usage
### Creating a Terminator

To create a new instance of the terminator, you need to specify the signals that should trigger the termination. The terminator listens for these signals and closes the registered resources when a signal is received.

```go

import (
	"os"
	"github.com/RohanPoojary/go-terminator"
)

func main() {
	closeSignals := []os.Signal{os.Interrupt, os.Kill}
	term := terminator.NewTerminator(closeSignals)
}
```

### Adding Resources

Resources that need to be closed gracefully can be registered with the terminator using the Add and AddWithTimeout methods. These methods take the resource name, a closing function, and an optional timeout duration.

```go

term.Add("Database Connection", func(ctx context.Context) error {
	// Close the database connection gracefully.
	return db.Close()
})

term.AddWithTimeout("File Writer", func(ctx context.Context) error {
	// Close the file writer gracefully, allowing a maximum of 5 seconds for closure.
	return fileWriter.Close()
}, 5*time.Second)
```

### Setting Callback

You can set a callback function that will be executed after all registered resources are closed. This can be useful for performing any final tasks or logging.

```go

term.SetCallback(func(result terminator.TerminationResult) {
	fmt.Println("Termination completed. Result:", result)
})
```

### Waiting for Termination

The Wait method allows you to wait for the termination process to complete with a specified timeout duration.

```go

success := term.Wait(10 * time.Second)
if success {
	fmt.Println("Termination completed successfully.")
} else {
	fmt.Println("Termination timed out.")
}
```


### TerminationResult Structure

The TerminationResult structure provides information about the termination process:

* `Signal`: The termination signal received.
* `Result`: A slice of TerminationResultData containing information about each closed resource.

## Contributing

Contributions are welcome! If you find any issues or have suggestions, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the Open Source Apache License - see the LICENSE file for details.