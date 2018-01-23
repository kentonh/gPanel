jQuery('._js_add-subdomain-form').on('submit', function(e){
    e.preventDefault();

    if(jQuery('#addSubdomain') && jQuery('#addSubdomain').val() && jQuery('#subdomainRoot') && jQuery('#subdomainRoot').val()) {
        var requestData = {};
        requestData["name"] = jQuery('#addSubdomain').val();
        requestData["root"] = jQuery('#subdomainRoot').val();

        var xhr = new XMLHttpRequest();
        xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function() {
            if(xhr.status == 204) {
                ListSubdomains();
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
    else {
        alert('All fields must be filled out to submit this form.');
    }
});