var ipModal = jQuery('.ip-filter-modal');

jQuery('._js_ip-filtering-open').on('click', function(e){
  e.preventDefault();

  var title;
  switch(jQuery(this).attr('data')) {
    case "general":
      title = "General";
      ipModal.find('input[name="type"]').attr('value', 'general');
      ipModal.find('#filterIPHelp').html("Filtering this IP under the general filter type will disallow access to the website for all modes.");
      break;
    case "maintenance":
      title = "Maintenance Mode";
      ipModal.find('input[name="type"]').attr('value', 'maintenance');
      ipModal.find('#filterIPHelp').html("Whitelisting this IP under the maintenance filter type will allow access to the website during maintenance mode.");
      break;
    default:
      alert("Error, refresh and try again. If problem persists contact server administrator.");
      return;
  }
  title += " IP Filtering";

  ipModal.find('.modal-title').html(title);

  listFilteredIPs(jQuery(this).attr('data'));

  ipModal.modal('show');
});

function listFilteredIPs(type) {
  ipModal.find('._js_currently-filtered-ips').html('');

  var requestData = {}
  requestData["type"] = type;

  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'api/ip/list', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      if(xhr.response != undefined && xhr.response.length != 0) {
        jsonResponse = JSON.parse(xhr.response)
        console.log(xhr.response);
        jQuery.each(jsonResponse, function(k, v) {
          ipModal.find('._js_currently-filtered-ips').append("<li>"+v.ip+"</li>");
        });
      }
      else {
        ipModal.find('.modal-body').html("An error has occurred, please refresh. If problem persists contact your administrator.");
      }
    }
    else if(xhr.status == 204) {
      ipModal.find('._js_currently-filtered-ips').append("<li>No Filtered IPs Currently Exist.</li>");
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        ipModal.find('.modal-body').html(xhr.response);
      }
      else {
        ipModal.find('.modal-body').html("An error has occurred, please refresh. If problem persists contact your administrator.");
      }
    }
  }
}
