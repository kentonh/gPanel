jQuery('._js_add-nameserver-form').on('submit', function(e){
    e.preventDefault();

    var nameserver = jQuery(this).find('#addNameserver');
    if(nameserver && nameserver.val()) {
        var requestData = {};
        requestData["nameserver"] = nameserver.val();

        var xhr = new XMLHttpRequest();
        xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
        xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function() {
            if (xhr.status == 204) {
                ListNameservers();
            }
            else {
                if(xhr.response != undefined && xhr.response.length != 0) {
                    alert('Error: ' + xhr.response);
                }
                else {
                    alert('An error has occurred, refresh and try again. If problem persists please contact your administrator.');
                }
            }
        }
    }
    else {
        alert('The nameserver field has to be filled out.');
    }
});