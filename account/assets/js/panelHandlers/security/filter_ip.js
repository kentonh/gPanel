var ipModal = jQuery('.ip-filter-modal');

jQuery('._js_ip-filter-form').on('submit', function(e){
  e.preventDefault();

  var requestData = {};
  requestData["ip"] = jQuery(this).find('input[name="ip"]').val();
  requestData["type"] = jQuery(this).find('input[name="type"]').val();

  var xhr = new XMLHttpRequest();
  xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if (xhr.status == 204) {
      ipModal.find('._js_currently-filtered-ips').append('<li>'+requestData["ip"]+'</li>');
    }
    else {
      alert("Something went wrong trying to filter that IP, please contact your administrator if problem persists.");
    }
  }
});
