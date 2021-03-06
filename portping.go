package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
TODO:
i orginally made this to be as simple as possible but saw the potential
for a load testing tool right away. It might be worth investigating if
i can specify the amount of data to send and then an option to control
the rate.
*/

var help_msg string = `
verify a tcp port is open on a remote machine

usage:
portping HOST PORT [count, timeout]

arguments:
-c    how many times to ping (default is forever)
-t    timeout, how long to wait before we consider it failed
-s    dont display stats (returns a 1 or 0 exit status)
-m    if running silently, you may want to see over all stats
-r    how long to wait between pings in ms, 0 is none. default is 500
`

var (
	host    string
	port    string
	count   string
	timeout string
	rate    string
	current int
	silent  bool
	metrics bool
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	raw_args := os.Args
	if len(raw_args) < 3 {
		fmt.Printf("error expetected at leat two args, HOST PORT\n%v", help_msg)
		os.Exit(1)
	} else {
		host = raw_args[1]
		port = raw_args[2]
	}
	// set defaults
	count = "1000000"
	timeout = "5"
	silent = false
	rate = "500"
	for i, a := range raw_args[1:] {
		if !strings.HasPrefix(a, "-") {
			continue
		} else if a == "-h" {
			fmt.Println(help_msg)
			os.Exit(0)
		} else if a == "-c" {
			count = raw_args[i+2]
		} else if a == "-t" {
			timeout = raw_args[i+2]
		} else if a == "-s" {
			silent = true
		} else if a == "-m" {
			metrics = true
		} else if a == "-r" {
			rate = raw_args[i+2]
		} else {
			fmt.Println(help_msg)
			fmt.Printf("unexpected argument: %v\n", a)
			os.Exit(1)
		}
	}
	c, err := strconv.Atoi(count)
	check(err)
	t, err := strconv.Atoi(timeout)
	check(err)
	r, err := strconv.Atoi(rate)
	times := make([]float64, 0, c)
	start := time.Now()
	if !silent {
		fmt.Println(start)
	}
	target := host + ":" + port
	for current = 0; current < c; current++ {
		start_time := time.Now()
		conn, err := net.DialTimeout("tcp", target, time.Duration(t)*time.Second)
		if err != nil {
			fmt.Printf("encountered an error when trying port %v on %v:\n%v\n", port, host, err)
			os.Exit(1)
		}
		defer conn.Close()
		end_time := time.Now()
		elapsed_time := float64(end_time.Sub(start_time)) / float64(time.Millisecond)

		if !silent {
			fmt.Printf("pinged %v on port %v, took %f milliseconds\n", host, port, elapsed_time)
		}
		if r > 0 {
			times = append(times, elapsed_time)
			time.Sleep(time.Duration(r) * time.Millisecond)
		}
	}
	end := time.Now()
	runtime := float64(end.Sub(start)) / float64(time.Millisecond)
	avg := runtime / float64(current)
	if r > 0 {
		var total float64 = 0
		for _, val := range times {
			total += val
		}
		runtime = total
		avg = total / float64(len(times))
	}
	if !silent || metrics {
		fmt.Printf("connected %v times over %v milliseconds, average is %v\n", current, runtime, avg)
	}
	os.Exit(0)
}
