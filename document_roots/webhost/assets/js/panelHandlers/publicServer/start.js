jQuery('._js_public-server-start').on('click', function(e){
  e.preventDefault();

  var xhr = new XMLHttpRequest();
  xhr.open('UPDATE', 'api/server/start', true);
  xhr.send();

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
