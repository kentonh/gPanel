jQuery('._js_link-domain-form').on('submit', function(e){
    e.preventDefault();

    if(jQuery('#linkDomain') && jQuery('#linkDomain').val()) {
        var requestData = {};
        requestData["domain"] = jQuery('#linkDomain').val();
        requestData["name"] = BUNDLE_NAME;

        var xhr = new XMLHttpRequest();
        xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
        xhr.send(JSON.stringify(requestData));

        xhr.onloadend = function() {
            if(xhr.status == 204) {
                ListDomains(BUNDLE_NAME);
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
        alert('The domain field must be filled out to submit this form.');
    }
});