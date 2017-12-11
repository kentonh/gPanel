jQuery(document).on('click', '._js_remove-nameserver', function(e){
    e.preventDefault();

    var nameserver = jQuery(this).attr('data');
    var requestData = {};
    requestData["nameserver"] = nameserver;

    var ensure = confirm("Are you sure you want to delete this nameserver?");

    if(ensure) {
        var xhr = new XMLHttpRequest();
        xhr.open('DELETE', 'api/settings/remove_nameserver', true);
        xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function () {
            if (xhr.status == 204) {
                ListNameservers();
            }
            else {
                if (xhr.response != undefined && xhr.response.length != 0) {
                    alert('Error: ' + xhr.response);
                }
                else {
                    alert('An error has occurred, refresh and try again. If problem persists please contact your administrator.');
                }
            }
        }
    }
});