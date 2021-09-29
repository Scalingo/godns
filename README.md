# GODNS

A simple and fast DNS cache server written in Go.


Similar to [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html), but supports some difference features:


* Keep hosts records in Redis and the local file /etc/hosts

* Auto-Reloads when hosts configuration is changed. (Yes, dnsmasq needs to be reloaded)


## Local Deployment

Start GoDNS:

```shell
docker-compose up
```

You may see the error:

```text
Error starting userland proxy: listen udp 0.0.0.0:5321: bind: address already in use
```

It means there is another process on your server listening on the UDP port 5321.
You can list the process listening on this port with:

```bash
sudo netstat -tulpn | grep 5321
```

DNS query to GoDNS:

```shell
dig @172.17.0.1 -p 5321 www.github.com
```

## Use GoDNS

		$ sudo vi /etc/resolv.conf
		nameserver #the ip of godns running

## Configuration

All the configuration in `godns.conf` is a TOML format config file.

#### resolv.conf

Upstream server can be configured by changing file from somewhere other than "/etc/resolv.conf"

```
[resolv]
resolv-file = "/etc/resolv.conf"
```
If multiple `namerservers` are set in resolv.conf, the upsteam server will try in a top to bottom order


#### server-list-file
Domain-specific nameservers configuration, formatting keep compatible with Dnsmasq.
>server=/google.com/8.8.8.8

More cases please refererence [dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list)


#### cache

Only the local memory storage backend is currently implemented.  The redis backend is in the todo list

```
[cache]
backend = "memory"
expire = 600  # default expire time 10 minutes
maxcount = 100000
```



#### hosts

Force resolve domain to assigned ip, support two types hosts configuration:

* locale hosts file
* remote redis hosts

__hosts file__

can be assigned at godns.conf,default : `/etc/hosts`

```
[hosts]
host-file = "/etc/hosts"
```
Hosts file format is described in [linux man pages](http://man7.org/linux/man-pages/man5/hosts.5.html).
More than that , `*.` wildcard is supported additional.


__redis hosts__

This is a special requirment in our system. Must maintain a global hosts configuration,
and support update the host records from other remote server.
Therefore, while "redis-hosts" be enabled, will query the redis db when each dns request is reached.

The hosts record is organized with redis hash map. and the key of the map is configured.

```
[hosts]
redis-key = "godns:hosts"
```

_Insert hosts records into redis_

```
redis > hset godns:hosts www.test.com 1.1.1.1
```

Compared with file-backend records, redis-backend hosts support multiple A entries.

```
redis > hset godns:hosts www.test.com 1.1.1.1,2.2.2.2
```

__Zone configuration__

Specific configuration if you want your godns instance to only respond to one DNS Zone

```
zone = 'example.test.'
zone-ns = "ns-cloud-b1.googledomains.com."
zone-mbox = "cloud-dns-hostmaster.google.com."
zone-serial = 1
zone-refresh = 21600
zone-retry = 3600
zone-expire = 259200
zone-negcache-ttl = 30
zone-soa-ttl = 3600
```

If `zone` is set, all the other parameters should be set. These parameters are
used when one domain is not found in the hosts (file or redis), we return a
"NODATA" dns response.  According to the RFC, it should contain the SOA entry
in order to be cached correctly by other DNS servers.

If another domain than `zone` is queried, GoDNS will respond with rcode `REFUSED`
since only the specified zone will be responded to.


## Benchmark


__Debug close__

```
$ go test -bench=.

testing: warning: no tests to run
PASS
BenchmarkDig-8     50000             57945 ns/op
ok      _/usr/home/keqiang/godns        3.259s
```

The result : 15342 queries/per second

The test environment:

CentOS release 6.4

* CPU:
Intel Xeon 2.40GHZ
4 cores

* MEM:
46G


## Web console

Joke: A web console for godns

[https://github.com/kenshinx/joke](https://github.com/kenshinx/joke)

screenshot

![joke](https://raw.github.com/kenshinx/joke/master/screenshot/joke.png)



## Deployment

Deployment in productive supervisord highly recommended.

```

[program:godns]
command=/usr/local/bin/godns -c /etc/godns.conf
autostart=true
autorestart=true
user=root
stdout_logfile_maxbytes = 50MB
stdoiut_logfile_backups = 20
stdout_logfile = /var/log/godns.log

```


## TODO

* The redis cache backend
* Update ttl

## LICENSE
godns is under the MIT license. See the LICENSE file for details.
