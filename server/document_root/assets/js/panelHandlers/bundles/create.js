newBundleModal = jQuery('.new-bundle-modal');

jQuery('._js_bundles-create').on('click', function(e){
  e.preventDefault();

  newBundleModal.modal('show');
});

jQuery('._js_create-bundle-form').on('submit', function(e){
  e.preventDefault();

  var formData = {};
  for(var y = 0, yy = this.length; y < yy; y++) {
    var input = this[y];
    if(input.name) {
      if(input.type == "number") {
        formData[input.name] = parseInt(input.value);
      }
      else {
        formData[input.name] = input.value;
      }
    }
  }

  var xhr = new XMLHttpRequest();
  xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
  xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
  xhr.send(JSON.stringify(formData));

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      if(xhr.response != undefined && xhr.response.length != 0) {
        alert("Bundle \"" + xhr.response + "\" successfully created.");
      }
      else {
        alert("Bundle successfully created.");
      }
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        alert("Error: " + xhr.response);
      }
      else {
        alert("An error has occurred. Please try again. If problem persists contact server administrator.");
      }
    }
  }
});
