var subdomainModal = jQuery('.subdomain-management-modal');

jQuery('._js_subdomain-management').on('click', function(e){
    e.preventDefault();

    ListSubdomains();
    subdomainModal.modal('show');
});

function ListSubdomains() {
    var list = jQuery('._js_registered-subdomains');

    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/subdomain/list', true);
    xhr.send();

    xhr.onloadend = function() {
        list.html('');
        if(xhr.status == 200) {
            jsonResponse = JSON.parse(xhr.response)
            jQuery.each(jsonResponse, function(k, v) {
                list.append('<div class="row mt-2"><div class="col-6"><p class="mb-0">'+k+'</p><small>Root: '+v.root+'</small></div><div class="col-6 d-flex justify-content-end"><div class="btn-group" role="group"><button class="btn btn-outline-danger _js_delete-registered-subdomain" data="'+k+'">Delete</button></div></div></div>');
            });
        }
        else if(xhr.status == 204) {
            list.html('<div class="row mt-2"><div class="col-6 d-flex align-items-center"><p>No registered subdomains exist for this account.</p></div></div>');
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