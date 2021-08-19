# portping
ping tcp ports with golang

```bash
verify a tcp port is open on a remote machine

usage:
portping HOST PORT [-c, -m, -r, -s, -t]

arguments:
-c    how many times to ping (default is forever)
-t    timeout, how long to wait before we consider it failed
-s    dont display stats (returns a 1 or 0 exit status)
-m    if running silently, you may want to see over all stats
-r    how long to wait between pings in ms, 0 is none. default is 500
```
<br>
**example**
<br>

```
portping your_host your_port -c 40 -s -m -r 80
connected 40 times over 1716.107469 milliseconds, average is 42.902686725
```
