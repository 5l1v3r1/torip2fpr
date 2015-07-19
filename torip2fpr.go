// Given Tor relay IP addresses, find their corresponding fingerprints.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	tor "git.torproject.org/user/phw/zoossh.git"
)

// Maps IP addresses to group IDs.
type AddrLookup map[string]string

var addrLookup AddrLookup = make(AddrLookup)

func loadAddresses(fileName string) {

	fd, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()

		separated := strings.SplitN(line, ":", 2)
		identifier := separated[0]
		ipAddrs := strings.Split(separated[1], ",")

		for _, ipAddr := range ipAddrs {
			ipAddr = strings.TrimSpace(ipAddr)
			log.Printf("Adding IP address %s to group ID %s.\n", ipAddr, identifier)
			addrLookup[ipAddr] = identifier
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseConsensus(channel chan string) {

	for path := range channel {
		log.Println(path)
		consensus, err := tor.ParseConsensusFile(path)
		if err != nil {
			log.Fatal(err)
		}

		for fingerprint, getStatus := range consensus.RouterStatuses {
			status := getStatus()
			identifier, ok := addrLookup[status.Address.String()]
			if ok {
				fmt.Printf("%s, %s, %s\n", fingerprint, status.Address.String(), identifier)
			}
		}
	}
}

func runExtraction(dataDir, ipAddrFile string) {

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	log.Printf("Trying to make use of %d CPUs.\n", numCPU)

	loadAddresses(ipAddrFile)

	chans := make([]chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		log.Printf("Starting go routine #%d.\n", i)
		chans[i] = make(chan string)
		go parseConsensus(chans[i])
	}

	idx := 0
	walkFiles := func(path string, info os.FileInfo, err error) error {
		if _, err := os.Stat(path); err != nil {
			log.Fatalf("File %s does not exist.\n", path)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		chans[idx] <- path
		idx = (idx + 1) % numCPU
		return nil
	}
	filepath.Walk(dataDir, walkFiles)

	for i := 0; i < numCPU; i++ {
		log.Printf("Closing channel %d.\n", i)
		close(chans[i])
	}
}

func main() {

	dataDir := flag.String("datadir", "", "Directory containing consensuses to analyse.")
	ipAddrFile := flag.String("addrfile", "", "File containing IP addresses.")

	flag.Parse()

	if *dataDir == "" {
		log.Fatal("No data directory given.  Please use the -datadir switch.")
	}

	if *ipAddrFile == "" {
		log.Fatal("No IP address file given.  Please use the -addrfile switch.")
	}

	runExtraction(*dataDir, *ipAddrFile)
}
