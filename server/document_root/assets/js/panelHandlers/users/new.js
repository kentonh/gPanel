var userModal = jQuery('.user-management-modal');

var usernameInput = jQuery('#addUserUsername');
var passwordInput = jQuery('#addUserPassword');
var passwordInputRetype = jQuery('#addUserPasswordRetype');

jQuery('._js_add-user-form').on('submit', function(e){
  e.preventDefault();

  if((usernameInput && usernameInput.val()) && (passwordInput && passwordInput.val()) && (passwordInputRetype && passwordInputRetype.val())) {
    if(passwordInput.val() == passwordInputRetype.val()) {
      var requestData = {};
      requestData["user"] = usernameInput.val();
      requestData["pass"] = passwordInput.val();

      var xhr = new XMLHttpRequest();
      xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
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
            alert('An error has occurred, refresh and try again. If problem persists please contact your administrator.');
          }
        }
      }
    }
    else {
      alert('Password fields do not match.');
    }
  }
  else {
    alert('All fields must contain values.');
  }
});

jQuery('._js_add-user-generate-password').on('click', function(e){
  e.preventDefault();
  var gen = "";
  var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-[]:;<>?";

  for (var i = 0; i < 32; i++) {
    gen += chars.charAt(Math.floor(Math.random() * chars.length));
  }

  toggleShowPassword(true);
  passwordInput.prop('value', gen);
  passwordInputRetype.prop('value', gen);
});

jQuery('._js_user-management-show-password').on('change', function(e){
  e.preventDefault();

  if(this.checked) {
    toggleShowPassword(true);
  }
  else {
    toggleShowPassword(false);
  }
});

function toggleShowPassword(show) {
  if(show) {
    jQuery('._js_user-management-show-password').prop('checked', true);
    passwordInput.attr('type', 'text');
    passwordInputRetype.attr('type', 'text');
  }
  else {
    jQuery('._js_user-management-show-password').prop('checked', false);
    passwordInput.attr('type', 'password');
    passwordInputRetype.attr('type', 'password');
  }
}
