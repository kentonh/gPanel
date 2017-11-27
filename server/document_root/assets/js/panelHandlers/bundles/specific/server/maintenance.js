jQuery('._js_specific-bundle-public-maintenance-mode').on('click', function(e){
  e.preventDefault();

  var name = jQuery('.specific-bundle-modal').attr('data');
  var requestData = {};
  requestData["bundle_name"] = name;

  var xhr = new XMLHttpRequest();
  xhr.open('UPDATE', 'api/server/maintenance', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 204) {
      getPublicServerStatus(name);
    }
    else {
      alert('An error has occurred, please contact your webhost administrator.');
    }
  }
});
