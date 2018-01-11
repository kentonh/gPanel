function PopulateKeyList() {
  var display = jQuery('._js_current-authorized-keys');
  display.html('');
  var requestData = {};
  requestData["username"] = BUNDLE_NAME;

  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'api/ssh/getkeys', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      if(xhr.response != undefined && xhr.response.length != 0) {
        jsonResponse = JSON.parse(xhr.response)
        jQuery.each(jsonResponse, function(k, v) {
          display.append('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><textarea cols="50" rows="7" readonly>'+v+'</textarea></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-danger _js_delete-ssh-key" data="'+v+'">Delete Key</button></div></div></div>');
        });
      }
      else {
        display.html('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p>An error has occurred, please refresh. If problem persists please contact your administrator.</p></div></div>');
      }
    }
    else if(xhr.status == 204) {
      display.html('<p>No authorized keys exist.</p>');
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        display.html('<p>Error: ' + xhr.response + '</p>');
      }
      else {
        display.html('<p>An error has occurred, please refresh. If problem persists please contact your administrator.</p>');
      }
    }
  }
}
