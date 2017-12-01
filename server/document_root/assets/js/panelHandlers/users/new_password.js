var userModal = jQuery('.user-management-modal');
var newPassModal = jQuery('.new-pass-modal');

var newPassword = jQuery('#updatePassword');
var newPasswordRetype = jQuery('#updatePasswordRetype');
var newPasswordUsername = jQuery('#updatePasswordUsername');

jQuery(document).on('click', '._js_user-management-new-password', function(e){
  e.preventDefault();

  if(!jQuery(this).attr('data') || jQuery(this).attr('data') == "") {
    alert("An error has occurred, please refresh and try again. If problem persists please contact your administrator.");
    return;
  }

  var username = jQuery(this).attr('data');
  newPasswordUsername.attr('value', username);

  newPassModal.find('.modal-title').html('Changing password for "'+username+'"');
  toggleShowPasswordNewPassword(false);

  userModal.modal('hide');
  newPassModal.modal('show');
});

jQuery('._js_back-to-user-management').on('click', function(e){
  e.preventDefault();

  newPassModal.modal('hide');
  userModal.modal('show');
});

jQuery('._js_update-password-form').on('submit', function(e){
  e.preventDefault();

  if((newPassword && newPassword.val()) && (newPasswordRetype && newPasswordRetype.val()) && (newPasswordUsername && newPasswordUsername.val())) {
    if(newPassword.val() == newPasswordRetype.val()) {
      var ensure = confirm("Are you sure you want to change the password of user \"" + newPasswordUsername.val() + "\"?");
      if(ensure) {
        var requestData = {};
        requestData["user"] = newPasswordUsername.val();
        requestData["pass"] = newPassword.val();

        var xhr = new XMLHttpRequest();
        xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
        xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function() {
          if(xhr.status == 204) {
            alert("Password successfully updated.");
            newPassModal.modal('hide');
            userModal.modal('show');
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
    }
    else {
      alert("Passwords must match.");
    }
  }
  else {
    alert("All fields need to be filled out");
  }
});

jQuery('._js_update-password-show-password').on('change', function(e){
  e.preventDefault();

  if(this.checked) {
    toggleShowPasswordNewPassword(true);
  }
  else {
    toggleShowPasswordNewPassword(false);
  }
});

jQuery('._js_update-password-generate-password').on('click', function(e){
  e.preventDefault();

  var genpass = generatePassword();

  toggleShowPasswordNewPassword(true);
  newPassword.prop('value', genpass);
  newPasswordRetype.prop('value', genpass);
});

function toggleShowPasswordNewPassword(show) {
  if(show) {
    jQuery('._js_update-password-show-password').prop('checked', true);
    newPassword.attr('type', 'text');
    newPasswordRetype.attr('type', 'text');
  }
  else {
    jQuery('._js_update-password-show-password').prop('checked', false);
    newPassword.attr('type', 'password');
    newPasswordRetype.attr('type', 'password');
  }
}
