var logModal = jQuery('.view-log-modal');

jQuery('._js_diagnostics-view-log').on('click', function(e){
  e.preventDefault();

  var logName = jQuery(this).attr('data');

  var title;
  switch(logName) {
    case "client_errors":
      title = "Client Error Log (4xx)"
      break;
    case "server_errors":
      title = "Server Error Log (5xx)"
      break;
    case "load_time":
      title = "Load Time Log"
      break;
    default:
      return;
      break;
  }
  jQuery(logModal).find('.modal-title').html(title);

  var requestData = {};
  requestData["name"] = logName+".log";

  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'api/logs/read', true);
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
      if(xhr.response) {
        alert(xhr.response);
      }
      else {
        alert('An error has occurred, please contact your webhost administrator.');
      }
    }
  }

});
