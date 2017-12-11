var domainManagementModal = jQuery('.domain-management-modal');

jQuery('._js_domain-management').on('click', function(e){
    e.preventDefault();

    ListDomains(BUNDLE_NAME);
    domainManagementModal.modal('show');
});

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