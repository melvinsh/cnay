package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/melvinsh/cnay/pkg/resolver"
)

func main() {
	debug := flag.Bool("d", false, "Enable debug output")
	listFile := flag.String("l", "", "Path to the file containing the list of hostnames")
	showHostname := flag.Bool("r", false, "Show original hostname in brackets")
	useProgressBar := flag.Bool("pb", false, "Enable progress bar")
	flag.Parse()

	reader, err := resolver.GetInputReader(*listFile, *debug)
	if err != nil {
		switch err.(type) {
		case *resolver.MissingInputError:
			displayManual()
		default:
			fmt.Printf("error: %s\n", err)
		}
		os.Exit(1)
	}

	hostnames := resolver.ReadHostnames(reader, *debug)

	uniqueIPs := resolver.ResolveHostnames(hostnames, *debug, *showHostname, *useProgressBar)

	for _, ip := range uniqueIPs {
		fmt.Println(ip)
	}
}

func displayManual() {
	fmt.Println("Usage:")
	fmt.Println("  Resolve hostnames to IP addresses based on input from STDIN or a file.")
	fmt.Println("  Only return unique and sorted IP addresses if the hostname has an A record and not a CNAME record.")
	fmt.Println("  If the CNAME record is on the same domain as the input, return the IP address.")
	fmt.Println("  When the -r flag is specified, display the original hostname in brackets.")
	fmt.Println()
	fmt.Println("  -l string")
	fmt.Println("        Path to the file containing the list of hostnames")
	fmt.Println("  -d")
	fmt.Println("        Enable debug output")
	fmt.Println("  -r")
	fmt.Println("        Show original hostname in brackets")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./resolver -l hostnames.txt")
	fmt.Println("  echo 'www.example.com' | ./resolver")
	fmt.Println("  ./resolver -r -l hostnames.txt")
}
