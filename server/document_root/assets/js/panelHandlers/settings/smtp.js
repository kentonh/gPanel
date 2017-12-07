var smtpModal = jQuery('.smtp-settings-modal');

jQuery('._js_smtp-credentials').on('click', function(e){
  e.preventDefault();

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'api/settings/get_smtp', true);
  xhr.send();

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      var resp = JSON.parse(xhr.response);

      if(resp["type"] == "crammd5") {
        jQuery('#smtpType').val(resp["type"]).change();
      }

      jQuery('#smtpUsername').val(resp["username"]);
      jQuery('#smtpServer').val(resp["server"]);
      jQuery('#smtpPort').val(resp["port"]);
    }
    smtpModal.modal('show');
  }
});

jQuery('._js_smtp-settings-form').on('submit', function(e){
  e.preventDefault();

  var flag = false;
  jQuery(this).find('input').each(function(i){
    if(jQuery(this) && jQuery(this).val()) return true;
    else {
      flag = true;
      return false;
    }
  });

  if(flag) {
    alert('All inputs need to be filled out.');
    return;
  }

  var requestData = {};
  requestData["type"] = jQuery(this).find('#smtpType').val();
  requestData["username"] = jQuery(this).find('#smtpUsername').val();
  requestData["password"] = jQuery(this).find('#smtpPassword').val();
  requestData["server"] = jQuery(this).find('#smtpServer').val();
  requestData["port"] = parseInt(jQuery(this).find('#smtpPort').val());

  var xhr = new XMLHttpRequest();
  xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
  xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 204) {
      alert("New SMTP Settings Connect Successfully and are Saved.");
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
});
