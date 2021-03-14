package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"net/http"

	flag "github.com/spf13/pflag"
)

type WaiterOptions struct {
	directory, port string
}

// Directories to look for in CWD if no value is passed
var defaultDirs = []string{
	"public",
	"static",
}

// Returns the first directory in defaultDirs that exists in CWD or an error
func getDefaultDir() (string, error) {
	var directory string

	for _, dir := range defaultDirs {

		// err tells us whether dir exists
		_, err := os.Stat(dir)
		if !os.IsNotExist(err) {
			directory = dir
			break
		}
	}

	if directory == "" {
		return "", errors.New("cannot find any directories to serve")
	}

	return directory, nil
}

func printHelp() {
	fmt.Print("Waiter is a static website server to quickly prototype static sites.\n\n")
	fmt.Println("When run with no commands, Waiter will automatically look for a './public' and then a './static' directory, if './public' does not exist, and serve the contents of the directory on port 3000.")
	fmt.Println("The first argument passed to Waiter will tell Waiter which directory to serve. By using the '--port' or '-p' option you can choose which port to serve your directory on.")
}

func parseArgs() WaiterOptions {
	// --port or -p overrides default (3000) port
	portPtr := flag.StringP("port", "p", "3000", "The port to serve on.")
	helpPtr := flag.BoolP("help", "h", false, "The port to serve on.")
	flag.Parse()

	if *helpPtr {
		printHelp()
		os.Exit(0)
	}

	// Holds the directory to serve
	var dir string = flag.Arg(0) // Initialized to the first argument passed from CLI

	// If no dir is given from CLI, check for defaultDirs' existence in CWD
	if dir == "" {
		var err error

		dir, err = getDefaultDir()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Check if the passed dir exists
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			log.Fatal("cannot find the directory:\n	", dir)
		}
	}

	// log.Println("\nDir: "+dir, "Port: "+*portPtr, "Help: "+strconv.FormatBool(*helpPtr))
	opts := WaiterOptions{dir, *portPtr}
	return opts
}

func main() {
	var args WaiterOptions = parseArgs()

	// Setup server with CLI provided or default directory
	staticHandler := http.FileServer(http.Dir(args.directory))
	http.Handle("/", staticHandler)

	fmt.Printf("\n☕ Serving on http://localhost:%s ☕\n", args.port)

	// Start server on given port
	if err := http.ListenAndServe(":"+args.port, nil); err != nil {
		log.Fatal(err)
	}
}
