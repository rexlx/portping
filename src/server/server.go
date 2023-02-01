package main

import (
	"fmt"
	"net"
	"time"

	pb "github.com/rexlx/portping/src/proto"
)

type Server struct {
	pb.ReacherServer
}

type Pinger struct {
	Addr    string
	Port    int32
	Count   int32
	Timeout int32
	Wait    int32
	Silent  bool
	Stats   []time.Duration
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
	// check if the value is NaN https://go.dev/src/math/bits.go func IsNaN(f float64) (is bool)
	if res != res {
		return 0.0
	}
	return res

}

func (s *Server) Ping(in *pb.PingRequest, stream pb.Reacher_PingServer) error {
	if in.Count < 1 {
		in.Count = 1e4
	}
	// setting a really low wait is a recipe for disaster, or your client is probably trying to spam a port.
	// the default behavior is to wait no less than 26 milliseconds between pings
	if in.Wait < 26 {
		in.Wait = 26
	}
	ping := Pinger{
		Addr:    in.Address,
		Port:    in.Port,
		Count:   in.Count,
		Wait:    in.Wait,
		Timeout: 5,
		Stats:   []time.Duration{},
	}
	for i := 0; i <= int(ping.Count); i++ {
		start := time.Now()
		// attempt to dial
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ping.Addr, ping.Port), time.Duration(ping.Timeout)*time.Second)
		// if we cant, log and retry
		if err != nil {
			stream.Send(&pb.PingResponse{
				Error: true,
				Msg:   fmt.Sprintf("encountered an error when dialing port %v on %v:\n%v\n", ping.Port, ping.Addr, err),
			})
			time.Sleep(time.Duration(ping.Wait) * time.Millisecond)
			continue
		}
		// otherwise we were able to connect, wait to close it
		defer conn.Close()
		// get the duration of the round trip
		elapsed_time := time.Since(start)
		// send the good news
		stream.Send(&pb.PingResponse{
			Error:    false,
			Msg:      fmt.Sprintf("pinged %v on port %v, took %.3fms\n", ping.Addr, ping.Port, float64(elapsed_time.Microseconds())/1000),
			Duration: int64(elapsed_time),
		})
		// record the statistics
		ping.Stats = append(ping.Stats, elapsed_time)
		if ping.Wait > 0 {
			time.Sleep(time.Duration(ping.Wait) * time.Millisecond)
		}
	}
	// send statistics
	stream.Send(&pb.PingResponse{
		Error:    false,
		Msg:      fmt.Sprintf("average duration was %v", ping.getAverageConnectionTime()),
		Duration: int64(ping.getAverageConnectionTime()),
	})
	return nil
}
