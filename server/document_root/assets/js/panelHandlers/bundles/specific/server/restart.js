jQuery('._js_specific-bundle-public-restart-graceful, ._js_specific-bundle-public-restart-forceful').on('click', function(e){
  e.preventDefault();

  var name = jQuery('.specific-bundle-modal').attr('data');
  var requestData = {};
  requestData["bundle_name"] = name

  if(jQuery(this).hasClass('_js_specific-bundle-public-restart-graceful')) {
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
      getPublicServerStatus(name);
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
