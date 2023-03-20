package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func main() {
	listFile := flag.String("l", "", "Path to the file containing the list of hostnames")
	debug := flag.Bool("d", false, "Enable debug output")
	showHostname := flag.Bool("r", false, "Show original hostname in brackets")
	flag.Parse()

	reader, err := getInputReader(*listFile, *debug)
	if err != nil {
		if *debug {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}

	hostnames := readHostnames(reader, *debug)

	uniqueIPs := resolveHostnames(hostnames, *debug, *showHostname)
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

func getInputReader(listFile string, debug bool) (io.Reader, error) {
	if listFile == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			displayManual()
			return nil, fmt.Errorf("no input provided")
		}
		return os.Stdin, nil
	}

	file, err := os.Open(listFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	return file, nil
}

func readHostnames(reader io.Reader, debug bool) []string {
	scanner := bufio.NewScanner(reader)
	hostnames := []string{}
	for scanner.Scan() {
		hostnames = append(hostnames, scanner.Text())
	}
	if err := scanner.Err(); err != nil && debug {
		fmt.Printf("Error reading input: %s\n", err)
	}
	return hostnames
}

func resolveHostnames(hostnames []string, debug, showHostname bool) []string {
	type result struct {
		ip       string
		hostname string
		err      error
	}

	results := make(chan result)

	for _, hostname := range hostnames {
		go func(hostname string) {
			ips, err := net.LookupIP(hostname)
			if err != nil {
				if debug {
					fmt.Printf("Error resolving %s: %s\n", hostname, err)
				}
				results <- result{err: err}
				return
			}

			cname, err := getFinalCNAME(hostname)
			if err == nil && !sameDomain(hostname, cname) {
				if debug {
					fmt.Printf("%s is an alias for %s, skipping\n", hostname, cname)
				}
				results <- result{err: fmt.Errorf("%s is an alias for %s", hostname, cname)}
				return
			}

			for _, ip := range ips {
				if ip.To4() != nil {
					ipStr := ip.String()
					results <- result{ip: ipStr, hostname: hostname}
				}
			}
		}(hostname)
	}

	ipSet := make(map[string]struct{})
	ipHostnameMap := make(map[string]string)

	for range hostnames {
		res := <-results
		if res.err != nil {
			continue
		}

		ipSet[res.ip] = struct{}{}
		ipHostnameMap[res.ip] = res.hostname
	}

	uniqueIPs := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		if showHostname {
			uniqueIPs = append(uniqueIPs, fmt.Sprintf("%s [%s]", ip, ipHostnameMap[ip]))
		} else {
			uniqueIPs = append(uniqueIPs, ip)
		}
	}
	sort.Strings(uniqueIPs)
	return uniqueIPs
}

func sameDomain(hostname, cname string) bool {
	hostnameDomain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return false
	}

	cnameDomain, err := publicsuffix.EffectiveTLDPlusOne(cname)
	if err != nil {
		return false
	}

	return strings.EqualFold(hostnameDomain, cnameDomain)
}

func getFinalCNAME(hostname string) (string, error) {
	var finalCNAME string
	for {
		cname, err := net.LookupCNAME(hostname)
		if err != nil {
			return "", err
		}

		cname = strings.TrimSuffix(cname, ".")

		if cname == hostname {
			break
		}

		finalCNAME = cname
		hostname = cname
	}

	if finalCNAME == "" {
		return "", fmt.Errorf("no CNAME record found for %s", hostname)
	}

	return finalCNAME, nil
}
