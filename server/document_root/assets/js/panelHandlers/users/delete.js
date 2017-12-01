jQuery(document).on('click', '._js_user-management-delete', function(e){
  e.preventDefault();

  if(!jQuery(this).attr('data') || jQuery(this).attr('data') == "") {
    alert("An error has occurred, please refresh and try again. If problem persists please contact your administrator.");
    return;
  }

  var ensure = confirm('Are you sure you want to delete the user "' + jQuery(this).attr('data') + '"?');
  if(ensure) {
    var requestData = {};
    requestData["user"] = jQuery(this).attr('data');

    var xhr = new XMLHttpRequest();
    xhr.open('UPDATE', 'api/user/delete', true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
      if(xhr.status == 204) {
        listCurrentUsers();
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
});
