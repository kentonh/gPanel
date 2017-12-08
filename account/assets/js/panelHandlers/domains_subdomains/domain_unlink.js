jQuery(document).on('click', '._js_delete-registered-domain', function(e){
    e.preventDefault();

    var domain = jQuery(this).attr('data');
    var ensure = confirm("Are you sure you want to unlink the domain \""+ domain +"\" from your account?");

    if(ensure) {
        var requestData = {};
        requestData["domain"] = domain;

        var xhr = new XMLHttpRequest();
        xhr.open('DELETE', 'api/domain/unlink', true);
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function () {
            if (xhr.status == 204) {
                ListDomains(BUNDLE_NAME);
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