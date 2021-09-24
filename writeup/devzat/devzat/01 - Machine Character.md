# Plot
- Enumerate and find landing page, chat and pets
- Step one enumerate pets.devzat.htb
	- find .git folder
	- Git-dumper and static code analysis will lead to RCE
- Being patrick on the box you can traverse to catherine with influxdb
	- InfluxDB is running in Docker
	- It is bound locally to 8086
	- Patrick will get mail upon login about influxdb as a hint from root
	- Crafting the JWT token patrick can move to catherine
- Catherine can read her mails (var/mail/catherine) and see that patrick installed a locally bound dev instance of devzat to try out *new feature*.
	- Dev instance is bound locally -> ssh localhost:8443
	- ./file command should let you send files content to chatroom
	- You will need the password from the mail
- Using the PW from mail you can use the new command and read roots id_rsa or root.txt by leveraging path traversal vulnerability


## Services to discover externally
- 22 ssh -> key auth only
- 80 apache (ip redirects to devzat.htb)
	- devzat.htb -> Landinpage (HTML5 up!)
	- pets.devzat.htb -> vulnerable go API (RCE)
		- also has .git with source
- 8000 secure version of devzat

## Services to discover after patrick on machine
- 8086 influxdb -> Get Catherines Creds
	- https://snyk.io/vuln/SNYK-GOLANG-GITHUBCOMINFLUXDATAINFLUXDBSERVICESHTTPD-1041719
- 8443 insecure version of devzat
	- source is @ /var/backups/devzat_dev.zip
	- together with secure main branch devzat_main.zip
	- you can diff the extracted folders to get a quick view of changes


# HTB Overview

| name | difficulty |
|---|---|
|devzat|medium|

**Requirements:**  
Taken from [https://app.hackthebox.eu/machines/submission/overview](https://app.hackthebox.eu/machines/submission/overview)
- 3 Steps typically
- Custom exploitation, but straight forward
- Path clear from context / hints, no rabbit holes
- Very easy binary exploitation and/or reverse engineering
- Generating simple scripts

