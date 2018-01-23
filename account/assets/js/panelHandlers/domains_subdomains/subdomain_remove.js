jQuery(document).on('click', '._js_delete-registered-subdomain', function(e){
    e.preventDefault();

    var subdomain = jQuery(this).attr('data');
    var ensure = confirm("Are you sure you want to delete the subdomain \""+ subdomain +"\" from your account?");

    if(ensure) {
        var requestData = {};
        requestData["name"] = subdomain;

        var xhr = new XMLHttpRequest();
        xhr.open('DELETE', 'api/subdomain/remove', true);
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function () {
            if (xhr.status == 204) {
                ListSubdomains();
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