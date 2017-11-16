var statusSpan = jQuery("._js_public-server-status");

function getPublicServerStatus() {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'api/server/status', true);
  xhr.send();

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      statusSpan.attr('class', '_js_public-server-status');

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
      statusSpan.attr('class', '_js_public-server-status').addClass('text-danger').html('ERROR');
    }
  }
}

// Run it once on load
getPublicServerStatus();
