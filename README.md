GODNS
====

A simple and fast dns cache server written by go.


Similar as [dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html) ,but support some difference features:


* Keep hosts records in redis instead of the local file /etc/hosts  

* Atuo-Reload when hosts configuration changed. (Yes,dnsmasq need restart)

* Cache records save in memory or redis configurable


## Install & Running

1. Install  

		$ go get github.com/kenshinx/godns


2. Build  

		$ cd $GOPATH/src/github.com/kenshinx/godns 
		$ go build -o godns *.go


3. Running  

		$ sudo ./godns -c godns.conf


4. Use

		$ sudo vi /etc/resolv.conf
		nameserver 127.0.0.1



## Configuration

All the configuration on `godns.conf` a TOML formating config file.   
More about Toml :[https://github.com/mojombo/toml](https://github.com/mojombo/toml)


#### resolv.conf

Upstream server can be configuration by change file from somewhere other that "/etc/resolv.conf"

```
[resolv]
resolv-file = "/etc/resolv.conf"
```
If multi `namerserver` set at resolv.conf, the upsteam server will try in order of up to botton



#### cache

Only the local memory storage backend implemented now.  The redis backend is in todo list

```
[cache]
backend = "memory"   
expire = 600  # default expire time 10 minutes
maxcount = 100000
```



#### hosts

Force resolv domain to assigned ip, support two types hosts configuration:

* locale hosts file
* remote redis hosts

__hosts file__  

can be assigned at godns.conf,default : `/etc/hosts`

```
[hosts]
host-file = "/etc/hosts"
```


__redis hosts__ 

This is a espeical requirment in our system. Must maintain a gloab hosts configuration, 
and support update the hosts record from other remote server.
so "redis-hosts" is be supported, and will query the reids when each dns request reached.  

The hosts record is organized with redis hash map. and the key of the map is configured.

```
[hosts]
redis-key = "godns:hosts"
```

_Insert hosts records into redis_

```
redis > hset godns:hosts www.sina.com.cn 1.1.1.1
```



## Benchmak


__Debug close__

```
$ go test -bench=.

testing: warning: no tests to run
PASS
BenchmarkDig-8     50000             57945 ns/op
ok      _/usr/home/keqiang/godns        3.259s
```

The result : 15342 queries/per second

The enviroment of test:

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





