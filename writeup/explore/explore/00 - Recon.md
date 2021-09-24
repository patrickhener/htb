# Nmap
```
Discovered open port 44099/tcp on 10.10.10.247
Discovered open port 42135/tcp on 10.10.10.247
Discovered open port 2222/tcp on 10.10.10.247
Discovered open port 59777/tcp on 10.10.10.247
```

# 2222
Banana Studio SSH

# 44099
```
> nc 10.10.10.247 44099

HTTP/1.0 400 Bad Request
Date: Thu, 29 Jul 2021 06:28:38 GMT
Content-Length: 22
Content-Type: text/plain; charset=US-ASCII
Connection: Close

Invalid request line: %    
```

# 42135

# 59777
```
> nc 10.10.10.247 59777

HTTP/1.0 400 Bad Request 
Content-Type: text/plain
Date: Thu, 29 Jul 2021 06:29:33 GMT

BAD REQUEST: Syntax error. Usage: GET /example/file.html%    
```

wget /sdcard/user.txt to get user
wget /sdcard/DCIM/creds.jpg gets ssh creds