# portping
ping tcp ports with golang

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
