var nameserverModal = jQuery('.nameserver-config-modal');

jQuery('._js_nameserver-config').on('click', function(e){
    e.preventDefault();

    ListNameservers();
    nameserverModal.modal('show');
});

function ListNameservers(){
    var display = jQuery('._js_current-nameservers');

    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/settings/get_nameservers', true);
    xhr.send();

    xhr.onloadend = function() {
        display.html('');

        if(xhr.status == 200) {
            if(xhr.response != undefined && xhr.response.length != 0) {
                jsonResponse = JSON.parse(xhr.response)
                jQuery.each(jsonResponse, function(k, v) {
                    display.append('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p class="mb-0">'+v+'</p></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-danger _js_remove-nameserver" data="'+v+'">Delete</button></div></div></div>');
                });
            }
            else {
                display.html('<p class="mt-2">An error has occurred, please refresh. If problem persists please contact your administrator.</p>');
            }
        }
        else if(xhr.status == 204) {
            display.html('<p class="mt-2">There are no nameservers stored in the server.</p>');
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