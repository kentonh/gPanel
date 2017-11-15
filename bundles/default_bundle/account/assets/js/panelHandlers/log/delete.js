var logModal = jQuery('.view-log-modal');

jQuery('._js_diagnostics-clear-log').on('click', function(e){
  e.preventDefault();

  var logName = jQuery(this).attr('data');

  var title;
  switch(logName) {
    case "client_errors":
      title = "Client Error Log";
      break;
    case "server_errors":
      title = "Server Error Log";
      break;
    case "load_time":
      title = "Load Time Log";
      break;
    default:
      return;
      break;
  }

  var ensure = confirm('Are you sure you want to clear ' + title + '?');
  if (ensure)  {
    var requestData = {};
    requestData["name"] = logName+".log";

    var xhr = new XMLHttpRequest();
    xhr.open('UPDATE', 'api/log/delete', true);
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
      if(xhr.status == 204) {
        jQuery(logModal).find('.modal-body textarea').html('The log is currently empty.');
      }
      else {
        if(xhr.response != undefined && xhr.response.length != 0) {
          alert(xhr.response);
        }
        else {
          alert("An error has occurred, please contact your webhost administrator.");
        }
      }
    }
  }
  else {
    return;
  }
});
