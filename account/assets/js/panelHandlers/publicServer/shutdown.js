jQuery('._js_public-server-shutdown-graceful, ._js_public-server-shutdown-forceful').on('click', function(e){
  e.preventDefault();
  var requestData = {};

  if(jQuery(this).hasClass('_js_public-server-shutdown-graceful')) {
    requestData["graceful"] = true;
  }
  else {
    requestData["graceful"] = false;
  }

  var xhr = new XMLHttpRequest();
  xhr.open('UPDATE', 'api/server/shutdown', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 204) {
      getPublicServerStatus();
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
