# portping
ping tcp ports with goalng

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
$ portping 192.168.1.87 22 -r 250 -c 10 -s -m
connected 10 times over 4.989891 milliseconds, average is 0.4989891
```
