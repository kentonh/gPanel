var statusSpan = jQuery("._js_specific-bundle-public-status");

function getPublicServerStatus(name) {
  var xhr = new XMLHttpRequest();

  var requestData = {};
  requestData["bundle_name"] = name;
  console.log(requestData);

  xhr.open('POST', 'api/server/status', true);
  xhr.send(JSON.stringify(requestData));

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      statusSpan.attr('class', '_js_specific-bundle-public-status');

      switch(parseInt(xhr.response)) {
        case 0:
          statusSpan.addClass('text-danger').html("OFF");
          break;
        case 1:
          statusSpan.addClass('text-success').html("ON");
          break;
        case 2:
          statusSpan.addClass('text-warning').html("MAINTENANCE");
          break;
        case 3:
          statusSpan.addClass('text-info').html("RESTARTING");
          break;
        default:
          console.log(xhr.response);
          statusSpan.addClass('text-danger').html('ERROR');
          break;
      }
    }
    else {
      console.log(xhr.response);
      statusSpan.attr('class', '_js_specific-bundle-public-status').addClass('text-danger').html('ERROR');
    }
  }
}
