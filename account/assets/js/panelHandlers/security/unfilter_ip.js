var ipModal = jQuery('.ip-filter-modal');

jQuery(document).on('click', '._js_delete-filtered-ip', function(e){
  e.preventDefault();

  var ip = jQuery(this).text();
  var id = jQuery(this).attr('data');

  var ensure = confirm('Are you sure you want to delete the IP filter for ' + ip + '?');
  if (ensure) {
    var requestData = {};
    requestData["id"] = parseInt(id);

    var xhr = new XMLHttpRequest();
    xhr.open('UPDATE', 'api/ip/unfilter', true);
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
      if(xhr.status == 204) {
        listFilteredIPs(ipModal.find('input[name="type"]').val());
      }
      else {
        if(xhr.response != undefined && xhr.response.length != 0) {
          alert("Error: " + xhr.response);
        }
        else {
          alert("An error has occurred, if the problem persists please contact your administrator.");
        }
      }
    }
  }
});
