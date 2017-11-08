jQuery('._js_public-server-restart-graceful, ._js_public-server-restart-forceful').on('click', function(e){
  e.preventDefault();
  var requestData = {};

  if(jQuery(this).hasClass('_js_public-server-restart-graceful')) {
    requestData["graceful"] = true;
  }
  else {
    requestData["graceful"] = false;
  }

  var xhr = new XMLHttpRequest();
  xhr.open('UPDATE', 'api/server/restart', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 204) {
      getPublicServerStatus();
      alert('Server successfully restarted.');
    }
    else {
      if(xhr.response) {
        alert(xhr.response);
      }
      else {
        alert('An error has occurred, please contact your webhost administrator.');
      }
    }
  }
});
