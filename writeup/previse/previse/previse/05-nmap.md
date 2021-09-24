# nmap

```bash
PORT   STATE SERVICE REASON         VERSION
22/tcp open  ssh     syn-ack ttl 63 OpenSSH 7.6p1 Ubuntu 4ubuntu0.3 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   2048 53:ed:44:40:11:6e:8b:da:69:85:79:c0:81:f2:3a:12 (RSA)
| ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDbdbnxQupSPdfuEywpVV7Wp3dHqctX3U+bBa/UyMNxMjkPO+rL5E6ZTAcnoaOJ7SK8Mx1xWik7t78Q0e16QHaz3vk2AgtklyB+KtlH4RWMBEaZVEAfqXRG43FrvYgZe7WitZINAo6kegUbBZVxbCIcUM779/q+i+gXtBJiEdOOfZCaUtB0m6MlwE2H2SeID06g3DC54/VSvwHigQgQ1b7CNgQOslbQ78FbhI+k9kT2gYslacuTwQhacntIh2XFo0YtfY+dySOmi3CXFrNlbUc2puFqtlvBm3TxjzRTxAImBdspggrqXHoOPYf2DBQUMslV9prdyI6kfz9jUFu2P1Dd
|   256 bc:54:20:ac:17:23:bb:50:20:f4:e1:6e:62:0f:01:b5 (ECDSA)
|_ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBCnDbkb4wzeF+aiHLOs5KNLPZhGOzgPwRSQ3VHK7vi4rH60g/RsecRusTkpq48Pln1iTYQt/turjw3lb0SfEK/4=
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

22 and 80
