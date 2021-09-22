---
title: 'Hack The Box - Writeup'
subtitle: 'Forge'
author: 'c1sc0'
date: \today{}
documentclass: scrartcl
pandoc-latex-color:
  - classes: [command]
    color: blue
titlepage: true
titlepage-text-color: "FFFFFF"
titlepage-color: "0c0d0e"
titlepage-rule-color: "8ac53e"
logo: "/home/patrick/htb/report/forge/images/badge.png"
logo-width: 250pt
header-left: 'Forge'
footer-left: 'c1sc0'
footer-right: 'Page \thepage'
---

<!--- Latex foo for table of content -->
\renewcommand*\contentsname{Table of Content}
\pagebreak
\tableofcontents
\pagebreak
<!--- Latex foo for table of content ends-->

# Overview
| IP | Difficulty |
| --- | --- |
| 10.10.11.111 | Medium |

## Nmap
`sudo nmap -sC -sV -oA nmap/forge -vvv 10.10.11.111`

```bash
PORT   STATE    SERVICE REASON         VERSION
21/tcp filtered ftp     no-response
22/tcp open     ssh     syn-ack ttl 63 OpenSSH 8.2p1 Ubuntu 4ubuntu0.3 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   3072 4f:78:65:66:29:e4:87:6b:3c:cc:b4:3a:d2:57:20:ac (RSA)
| ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC2sK9Bs3bKpmIER8QElFzWVwM0V/pval09g7BOCYMOZihHpPeE4S2aCt0oe9/KHyALDgtRb3++WLuaI6tdYA1k4bhZU/0bPENKBp6ykWUsWieSSarmd0sfekrbcqob69pUJSxIVzLrzXbg4CWnnLh/UMLc3emGkXxjLOkR1APIZff3lXIDr8j2U3vDAwgbQINDinJaFTjDcXkOY57u4s2Si4XjJZnQVXuf8jGZxyyMKY/L/RYxRiZVhDGzEzEBxyLTgr5rHi3RF+mOtzn3s5oJvVSIZlh15h2qoJX1v7N/N5/7L1RR9rV3HZzDT+reKtdgUHEAKXRdfrff04hXy6aepQm+kb4zOJRiuzZSw6ml/N0ITJy/L6a88PJflpctPU4XKmVX5KxMasRKlRM4AMfzrcJaLgYYo1bVC9Ik+cCt7UjtvIwNZUcNMzFhxWFYFPhGVJ4HC0Cs2AuUC8T0LisZfysm61pLRUGP7ScPo5IJhwlMxncYgFzDrFRig3DlFQ0=
|   256 79:df:3a:f1:fe:87:4a:57:b0:fd:4e:d0:54:c6:28:d9 (ECDSA)
| ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH67/BaxpvT3XsefC62xfP5fvtcKxG2J2di6u8wupaiDIPxABb5/S1qecyoQJYGGJJOHyKlVdqgF1Odf2hAA69Y=
|   256 b0:58:11:40:6d:8c:bd:c5:72:aa:83:08:c5:51:fb:33 (ED25519)
|_ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILcTSbyCdqkw29aShdKmVhnudyA2B6g6ULjspAQpHLIC
80/tcp open     http    syn-ack ttl 63 Apache httpd 2.4.41 ((Ubuntu))
|_http-title: Did not follow redirect to http://forge.htb
|_http-server-header: Apache/2.4.41 (Ubuntu)
| http-methods: 
|_  Supported Methods: GET HEAD POST OPTIONS
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel
```

## /etc/hosts
80 redirects to `forge.htb`. So adding it to `/etc/hosts`

![[Pasted image 20210922093500.png]]

## Website

![[Pasted image 20210922093610.png]]

Interesting "Upload an image" button top right

![[Pasted image 20210922094015.png]]

Looks like you can either provide file or enter URL.

Uploading images works, whereas uploading a cmd shell for example doesn't.

If you try and choose to upload from URL the box will callback to you:

```bash
> sudo ncat -lnvp 80
[sudo] password for patrick: 
Ncat: Version 7.92 ( https://nmap.org/ncat )
Ncat: Listening on :::80
Ncat: Listening on 0.0.0.0:80
Ncat: Connection from 10.10.11.111.
Ncat: Connection from 10.10.11.111:38550.
GET /foo.png HTTP/1.1
Host: 10.10.14.8
User-Agent: python-requests/2.25.1
Accept-Encoding: gzip, deflate
Accept: */*
Connection: keep-alive
```

## Subdomain enumeration
Wfuzz will reveal another subdomain:

```bash
> wfuzz -c -w ~/tools/wordlists/SecLists/Discovery/DNS/subdomains-top1million-5000.txt -u 'http://forge.htb' -H "Host: FUZZ.forge.htb" --hw 26
********************************************************
* Wfuzz 3.1.0 - The Web Fuzzer                         *
********************************************************

Target: http://forge.htb/
Total requests: 4989

=====================================================================
ID           Response   Lines    Word       Chars       Payload                                                                                      
=====================================================================

000000024:   200        1 L      4 W        27 Ch       "admin"                                                                                      

Total time: 0
Processed Requests: 4989
Filtered Requests: 4988
Requests/sec.: 0
```

So adding it to `/etc/hosts` and again look at the resulting page.

## admin.forge.htb
![[Pasted image 20210922094923.png]]

So the idea is to leverage a vulnerablity at the upload from URL part to look at `admin.forge.htb` from internally.

![[Pasted image 20210922095028.png]]