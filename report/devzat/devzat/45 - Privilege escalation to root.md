# Broken Dev Branch of *devzat*
There is the locally running instance of *devzat* on Port 8443. Catherine will have a notification about a new mail when she logs in via ssh. The mail will give her instructions on how to use the instance and the new feature. The new feature suffers from a path traversal vulnerability and is running as root. Thus, catherine is able to see root's private key (id_rsa) or root.txt.

## Vulnerable function
The function including the file to print is vulnerable to a Path Traversal attack. Also it has a hard coded secret.

```go
func fileCommand(u *user, args []string) {
	[... snip ...]

	path := args[0]
	pass := args[1]

	// Check my secure password
	if pass != "CeilingCatStillAThingIn2021?" {
		u.system("You did provide the wrong password")
		return
	}

	// Get CWD
	cwd, err := os.Getwd()
	if err != nil {
		u.system(err.Error())
	}

	// Construct path to print
	printPath := filepath.Join(cwd, path)
	
	[... snip ...]
```
As one can see the user can provide the *path* of the file to be printed as first argument to the function. The secret is hardcoded and the user needs to provide this as second argument to the fucntion. The app will then user **filepath.Join** to build a path out of the **current working directory** and the user provided **path**. When joining user provided content is not sanitized at all. So providing `../../../../../../../../etc/passwd` as first argument to the function will get you the /etc/passwd file. You get the idea.

As the app is running as root the attacker could either directly print the root flag or print `/root/.ssh/id_rsa` and then ssh into the machine.

## Hints on how to progress
Either chat instance will give you hints on how to progress the box (in form of a forged chat history) if you login as `patrick`, `catherine` or `admin`. If you are not logging in from `127.0.0.1` you will be forced to change your nick.

## Zip archives
The archives with the source of both, the main branch and the dev branch, will be at `/var/backups/`. A player will be able to scp them to his machine and then unzipping them into 2 separate folders. Doing a diff will result in the display of the changes and a player can clearly see the hard coded secret and the vulnerable function, without being overloaded with source code.

A user could also just read the complete code or maybe compile and tinker with on his attackers box. All he would need to have ready is golang 1.16.X or newer.

## Autostart via systemd
```bash
[Unit]
Description=Devzat - Dev
After=network.target

[Service]
Type=simple
User=patrick
Group=patrick
Restart=always
RestartSec=5s

WorkingDirectory=/root/devzat-dev/
ExecStart=/root/devzat-dev/start.sh
SyslogIdentifier=devzat-dev

[Install]
WantedBy=multi-user.target
```