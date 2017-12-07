var adminSettingsModal = jQuery('.admin-settings-modal');

jQuery('._js_admin-settings').on('click', function(e){
    e.preventDefault();

    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/settings/get_admin', true);
    xhr.send();

    xhr.onloadend = function() {
        if(xhr.status == 200) {
            var resp = JSON.parse(xhr.response);

            jQuery('#adminName').val(resp["name"]);
            jQuery('#adminEmail').val(resp["email"]);
        }
        adminSettingsModal.modal('show');
    }
});

jQuery('._js_admin-settings-form').on('submit', function(e){
    e.preventDefault();

    var flag = false;
    jQuery(this).find('input').each(function(i){
        if(jQuery(this) && jQuery(this).val()) return true;
        else {
            flag = true;
            return false;
        }
    });

    if(flag) {
        alert('All inputs need to be filled out.');
        return;
    }

    var requestData = {};
    requestData["name"] = jQuery(this).find('#adminName').val();
    requestData["email"] = jQuery(this).find('#adminEmail').val();

    var xhr = new XMLHttpRequest();
    xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    xhr.send(JSON.stringify(requestData));

    xhr.onloadend = function() {
        if(xhr.status == 204) {
            alert('Administrator settings successfully set.');
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
});