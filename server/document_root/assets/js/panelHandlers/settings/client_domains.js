var clientDomainsModal = jQuery('.registered-client-domains-modal');

jQuery('._js_registered-domains').on('click', function(e){
    e.preventDefault();

    ListClientDomains();
    clientDomainsModal.modal('show');
});

jQuery(document).on('click', '._js_delete-client-registered-domain', function(e){
    e.preventDefault();

    var domain = jQuery(this).attr('data');
    var ensure = confirm("Are you sure you want to unregister the domain \""+ domain +"\" from the client?");

    if(ensure) {
        var requestData = {};
        requestData["domain"] = domain;

        var xhr = new XMLHttpRequest();
        xhr.open('DELETE', 'api/domain/unlink', true);
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function () {
            if (xhr.status == 204) {
                ListClientDomains();
            }
            else {
                if (xhr.response != undefined && xhr.response.length != 0) {
                    alert('Error: ' + xhr.status);
                }
                else {
                    alert('An error has occurred. If problem persists please contact your community administrator.');
                }
            }
        }
    }
});

function ListClientDomains() {
    var list = jQuery('._js_current-registered-client-domains');

    var requestData = {};
    requestData["name"] = "*";

    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'api/domain/list', true);
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
        list.html('');
        if(xhr.status == 200) {
            jsonResponse = JSON.parse(xhr.response)
            jQuery.each(jsonResponse, function(k, v) {
                list.append('<div class="row mt-2"><div class="col-6"><p class="mb-0">'+k+'</p><small>Bundle: '+v.name+'</small></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-danger _js_delete-client-registered-domain" data="'+k+'">Delete</button></div></div></div>');
            });
        }
        else if(xhr.status == 204) {
            list.html('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p>No client registered domains exist.</p></div></div>');
        }
        else {
            if(xhr.response != undefined && xhr.response.length != 0) {
                alert('Error: ' + xhr.status);
            }
            else {
                alert('An error has occurred. If problem persists please contact your community administrator.');
            }
        }
    }
}