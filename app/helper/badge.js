var fs = require('fs');
var page = require('webpage').create();

fs.changeWorkingDirectory('%%tmpdirhere%%');

page.viewportSize = { width: 660, height: 150 };
page.zoomFactor = 3;

page.open('badge.html', function () {
  var bb = page.evaluate(function () {
    return document
      .getElementsByClassName('wrapper')[0]
      .getBoundingClientRect();
  });

  page.clipRect = {
    top: bb.top * 3,
    left: bb.left * 3,
    width: bb.width * 3,
    height: bb.height * 3,
  };

  page.render('badge.png');
  phantom.exit();
});
