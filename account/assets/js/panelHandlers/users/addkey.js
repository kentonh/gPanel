var sshModal = jQuery(".manage-ssh-modal");

jQuery('._js_manage-ssh').on('click', function(e){
  e.preventDefault();

  PopulateKeyList();
  sshModal.modal('show');
});

jQuery('._js_add-key-form').on('submit', function(e){
  e.preventDefault();

  var key = jQuery('#addKey');
  if(key && key.val()) {
    var requestData = {};
    requestData["username"] = BUNDLE_NAME;
    requestData["publicKey"] = key.val();

    var xhr = new XMLHttpRequest();
    xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
      if(xhr.status == 204) {
        PopulateKeyList();
        alert("Successfully added authorized public key.");
      }
      else {
        if(xhr.response != undefined && xhr.response.length != 0) {
          alert('Error: ' + xhr.response);
        }
        else {
          alert('An error has occurred, refresh and try again. If problem persists please contact your administrator.');
        }
      }
    }
  }
  else {
    alert("All fields must be filled out.")
  }
});
