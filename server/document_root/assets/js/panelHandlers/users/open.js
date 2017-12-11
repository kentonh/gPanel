var userModal = jQuery('.user-management-modal');

jQuery('._js_manage-users').on('click', function(e){
  e.preventDefault();

  jQuery('._js_user-management-show-password').prop('checked', false);

  listCurrentUsers();
  userModal.modal('show');
});

function listCurrentUsers() {
  var display = jQuery('._js_current-users');
  display.html('');
  var requestData = {};

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'api/user/list', true);
  xhr.send();

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      if(xhr.response != undefined && xhr.response.length != 0) {
        jsonResponse = JSON.parse(xhr.response)
        jQuery.each(jsonResponse, function(k, v) {
          display.append('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p class="mb-0">'+v+'</p></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-primary _js_user-management-new-password" data="'+v+'">New Password</button><button class="btn btn-outline-danger _js_user-management-delete" data="'+v+'">Delete</button></div></div></div>');
        });
      }
      else {
        display.html('<p class="mt-2">An error has occurred, please refresh. If problem persists please contact your administrator.</p>');
      }
    }
    else if(xhr.status == 204) {
      display.html('<p class="mt-2">There are no users in the server. This is a problem, this shouldn\'t be like this.</p>');
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        display.html('<p class="mt-2">Error: ' + xhr.response + '</p>');
      }
      else {
        display.html('<p class="mt-2">An error has occurred, please refresh. If problem persists please contact your administrator.</p>');
      }
    }
  }
}
