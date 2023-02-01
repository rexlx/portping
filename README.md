# portping
ping tcp ports with go

```bash
confirm a specific port is open on a target machine. usage:
portping <address> <port> [args]
portping pandora.com 443 -c 10 -w 50 -s

-c  count,   how many times to ping the target
-h  help,    print this help message and exit
-s  silent,  suppress output
-t  timeout, how long to wait before failing  
-w  wait,    how long to wait between pings. default is 500ms
```

to run the grpc server you'll need the protoc tool and go compiler.
```
./protogen.sh
./server       # you may need to open a port up server side
./client       # this will need to be configured to point to the server

#example client output
rxlx@rxmb portping % ./client
2023/01/31 18:58:03 8721833	pinged 8.8.8.8 on port 53, took 8.721ms
2023/01/31 18:58:04 9451500	pinged 8.8.8.8 on port 53, took 9.451ms
2023/01/31 18:58:04 10264083	pinged 8.8.8.8 on port 53, took 10.264ms
2023/01/31 18:58:05 10145375	pinged 8.8.8.8 on port 53, took 10.145ms
2023/01/31 18:58:05 9079333	pinged 8.8.8.8 on port 53, took 9.079ms
2023/01/31 18:58:06 9224166	pinged 8.8.8.8 on port 53, took 9.224ms
2023/01/31 18:58:06 9	average duration was 9.481
```
where the third column is the response time in nanoseconds.

