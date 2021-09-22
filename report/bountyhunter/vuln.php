<?php

if(isset($_POST['data'])) {
$xml = base64_decode($_POST['data']);
libxml_disable_entity_loader(false);
$dom = new DOMDocument();
$dom->loadXML($xml, LIBXML_NOENT | LIBXML_DTDLOAD);
$bugreport = simplexml_import_dom($dom);
}
?>
If DB were ready, would have added:
<table>
  <tr>
    <td>Title:</td>
    <td><?php echo $bugreport->title; ?></td>
  </tr>
  <tr>
    <td>CWE:</td>
    <td><?php echo $bugreport->cwe; ?></td>
  </tr>
  <tr>
    <td>Score:</td>
    <td><?php echo $bugreport->cvss; ?></td>
  </tr>
  <tr>
    <td>Reward:</td>
    <td><?php echo $bugreport->reward; ?></td>
  </tr>
</table>
