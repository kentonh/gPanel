var domainManagementModal = jQuery('.domain-management-modal');

jQuery('._js_domain-management').on('click', function(e){
    e.preventDefault();

    ListDomains(BUNDLE_NAME);
    ListServerNameservers();
    domainManagementModal.modal('show');
});

function ListServerNameservers() {
    var display = jQuery('._js_server-nameservers');

    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/settings/get_nameservers', true);
    xhr.send();

    xhr.onloadend = function() {
        display.html('');

        if(xhr.status == 200) {
            if(xhr.response != undefined && xhr.response.length != 0) {
                jsonResponse = JSON.parse(xhr.response)
                jQuery.each(jsonResponse, function(k, v) {
                    display.append('<div class="row mt-2"><div class="col-12 d-flex align-items-center"><p class="mb-0">'+v+'</p></div></div>');
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

function ListDomains(bundle_name) {
    var list = jQuery('._js_registered-domains');

    var requestData = {};
    requestData["name"] = bundle_name;

    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'api/domain/list', true);
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
        list.html('');
        if(xhr.status == 200) {
            jsonResponse = JSON.parse(xhr.response)
            jQuery.each(jsonResponse, function(k, v) {
                list.append('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p class="mb-0">'+k+'</p></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-danger _js_delete-registered-domain" data="'+k+'">Delete</button></div></div></div>');
            });
        }
        else if(xhr.status == 204) {
            list.html('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p>No registered domains exist for this account.</p></div></div>');
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