<%
Response.ContentType = "application/x-netscape-revocation"
serialnumber = Request.QueryString
set Admin = Server.CreateObject("CertificateAuthority.Admin")

stat = Admin.IsValidCertificate("earth.windcorp.htb\windcorp-CA", serialnumber)

if stat = 3 then Response.Write("0") else Response.Write("1") end if
%>
