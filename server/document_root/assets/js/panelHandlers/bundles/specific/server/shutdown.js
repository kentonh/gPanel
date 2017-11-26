jQuery('._js_specific-bundle-public-shutdown-graceful, ._js_specific-bundle-public-shutdown-forceful').on('click', function(e){
  e.preventDefault();

  var name = jQuery('.specific-bundle-modal').attr('data');
  var requestData = {};
  requestData["bundle_name"] = name;

  if(jQuery(this).hasClass('_js_specific-bundle-public-shutdown-graceful')) {
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
      getPublicServerStatus(name);
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
