/* Store bundle name in global variable */
var BUNDLE_NAME = "";

var xhr = new XMLHttpRequest();
xhr.open('GET', 'api/settings/name', true);
xhr.send();

xhr.onloadend = function() {
    if(xhr.status == 200) {
        BUNDLE_NAME = xhr.response;
    }
    else {
        if (xhr.response != undefined && xhr.response.length != 0) {
            alert('Error getting bundle name: ' + xhr.status);
        }
        else {
            alert('An error has occurred while getting the bundle name. If problem persists please contact your community administrator.');
        }
    }
}