jQuery('._js_delete-specific-bundle').on('click', function(e){
  e.preventDefault();

  var bundle = jQuery('.specific-bundle-modal').attr('data');

  var ensure = confirm('Are you sure you want to delete this bundle? This action CANNOT be undone.');
  if(ensure) {
    var requestData = {};
    requestData["name"] = bundle;
    requestData["bundle_name"] = bundle;

    var xhr = new XMLHttpRequest();
    xhr.open('DELETE', 'api/bundle/delete', true);
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
      if(xhr.status == 204) {
        alert('Bundle \"' + bundle + '\" has been successfully deleted.');
        jQuery('.specific-bundle-modal').modal('hide');
      }
      else {
        if(xhr.response != undefined && xhr.response.length != 0) {
          alert('Error: ' + xhr.response);
        }
        else {
          alert('An error has occurred. Please refresh and try again. If error persists please contact your administrator.');
        }
      }
    }
  }
});
