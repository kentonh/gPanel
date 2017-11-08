jQuery('._js_public-server-maintenance-mode').on('click', function(e){
  e.preventDefault();

  var xhr = new XMLHttpRequest();
  xhr.open('UPDATE', 'api/server/maintenance', true);
  xhr.send();

  xhr.onloadend = function() {
    if(xhr.status == 204) {
      getPublicServerStatus();
    }
    else {
      alert('An error has occurred, please contact your webhost administrator.');
    }
  }
});
