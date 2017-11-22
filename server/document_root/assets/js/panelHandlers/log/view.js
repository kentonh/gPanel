var logModal = jQuery('.view-log-modal');

jQuery('._js_log-view').on('click', function(e){
  e.preventDefault();

  var logName = jQuery(this).attr('data');

  var title;
  switch(logName) {
    case "server_errors":
      title = "Server Error Log";
      break;
    default:
      return;
      break;
  }
  jQuery(logModal).find('.modal-title').html(title);
  jQuery(logModal).find('._js_log-clear').attr('data', logName);

  var requestData = {};
  requestData["name"] = logName;

  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'api/log/read', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      var responseData;

      if(xhr.response != undefined && xhr.response.length != 0) {
        responseData = xhr.response;
      }
      else {
        responseData = "The log is currently empty."
      }

      jQuery(logModal).find('.modal-body textarea').html(responseData);
      jQuery(logModal).modal('show');
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        alert(xhr.response);
      }
      else {
        alert('An error has occurred, please contact your webhost administrator.');
      }
    }
  }

});
