jQuery(document).on('click', '._js_delete-ssh-key', function(e){
  e.preventDefault();

  var keyData = jQuery(this).attr('data');
  if(keyData != "") {
    var ensure = confirm('Are you sure you want to delete the specified key?');
    if(ensure) {
      var requestData = {};
      requestData["username"] = BUNDLE_NAME;
      requestData["publickey"] = keyData;

      var xhr = new XMLHttpRequest();
      xhr.open('UPDATE', 'api/ssh/deletekey', true);
      xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
      xhr.send(JSON.stringify(requestData));

      xhr.onloadend = function() {
        if(xhr.status == 204) {
          PopulateKeyList();
        }
        else {
          if(xhr.response != undefined && xhr.response.length != 0) {
            alert('Error: ' + xhr.response);
          }
          else {
            alert("An error has occurred, please refresh and try again. If problem persists please contact your administrator.");
          }
        }
      }
    }
  } else {
    alert("An error has occurred, please refresh and try again. If problem persists please contact your administrator.");
  }
});
