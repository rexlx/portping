package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var help string = `
confirm a specific port is open on a target machine. usage:
./portping <address> <port> [args]
./portping pandora.com 443 -c 10 -w 50 -s

-c  count,   how many times to ping the target
-h  help,    print this help message and exit
-s  silent,  only print end results
-t  timeout, how long to wait before failing  
-w  wait,    how long to wait between pings. default is 500ms
`

type Pinger struct {
	Addr    string
	Port    int
	Count   int
	Timeout int
	Wait    int
	Silent  bool
	Stats   []time.Duration
}

// state represents the exit condition and is passed to os.Exit
var state int

// parseArgs modifies the default pinger with the user supplied args
func (p *Pinger) parseArgs() {
	if len(os.Args) < 3 {
		log.Fatalln("expected two args, <host> <port>")
	}
	// we use positional args. the address will always go first
	p.Addr = os.Args[1]
	p.Port = strToInt(os.Args[2])
	// empty slice for our duration statistics
	p.Stats = []time.Duration{}
	// parse the args
	for i, a := range os.Args[1:] {
		switch a {
		case "-c":
			p.Count = strToInt(os.Args[i+2])
		case "-t":
			p.Timeout = strToInt(os.Args[i+2])
		case "-w":
			p.Wait = strToInt(os.Args[i+2])
		case "-s":
			p.Silent = true
		case "-h":
			log.Fatal(help)
		default:
		}
	}
}

// getAverageConnectionTime uses a slice of time durations and returns their average in ms
func (p *Pinger) getAverageConnectionTime() float64 {
	var t time.Duration
	// iter over the slice of durations and increment our artificial one
	for _, i := range p.Stats {
		t += i
	}
	// we technically have precision down to the nanosecond, but this is a reasonable program
	// for reasonable people
	res := (float64(t.Microseconds()) / 1000) / float64(len(p.Stats))
	// check if the value is NaN as done here: https://go.dev/src/math/bits.go func IsNaN(f float64) (is bool)
	if res != res {
		return 0.0
	}
	return res

}

// strToInt returns a converted integer or dies
func strToInt(s string) int {
	out, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func main() {
	// mark the start time
	begin := time.Now()

	// init a new pinger with some non-zero values
	ping := Pinger{
		Count:   1e4,
		Timeout: 5,
		Wait:    500,
	}
	// get the args
	ping.parseArgs()

	// begin the work
	for i := 0; i <= ping.Count; i++ {
		// we time the duration of each ping using this
		start := time.Now()
		// attempt to dial
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ping.Addr, ping.Port), time.Duration(ping.Timeout)*time.Second)
		// if we cant, log and exit
		if err != nil {
			// the program will exit with 1 in the event of ANY failure
			state = 1
			if !ping.Silent {
				fmt.Printf("encountered an error when dialing port %v on %v:\n%v\n", ping.Port, ping.Addr, err)
			}
			time.Sleep(time.Duration(ping.Wait) * time.Millisecond)
			// try again
			continue
		}
		// otherwise we were able to connect, wait to close
		defer conn.Close()
		// get the duration of the round trip
		elapsed_time := time.Since(start)

		if !ping.Silent {
			fmt.Printf("pinged %v on port %v, took %.3fms\n", ping.Addr, ping.Port, float64(elapsed_time.Microseconds())/1000)
		}
		// add the time it took to do everything into our stats slice
		ping.Stats = append(ping.Stats, elapsed_time)
		if ping.Wait > 0 {
			time.Sleep(time.Duration(ping.Wait) * time.Millisecond)
		}
	}
	if !ping.Silent {
		runtime := time.Since(begin)
		fmt.Printf("ran for %v seconds. average connection time was %.3fms\n", runtime.Seconds(), ping.getAverageConnectionTime())
	}
	os.Exit(state)
}
