manageBundlesModal = jQuery(".manage-bundles-modal");

jQuery('._js_bundles-manage').on('click', function(e){
  e.preventDefault();

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'api/bundle/list', true);
  xhr.send();

  xhr.onloadend = function() {
    if(xhr.status == 200) {
      if(xhr.response != undefined && xhr.response.length != 0) {
        manageBundlesModal.find('.modal-body').html(xhr.response)
      }
      else {
        manageBundlesModal.find('.modal-body').html("An error has occurred. Please try again. If problem persists contact server administrator.")
      }
      manageBundlesModal.modal('show');
    }
    else if(xhr.status == 204) {
      manageBundlesModal.modal('show');
    }
    else {
      if(xhr.response != undefined && xhr.response.length != 0) {
        manageBundlesModal.find('.modal-body').html(xhr.response)
      }
      else {
        manageBundlesModal.find('.modal-body').html(xhr.status + " Error!")
      }
      manageBundlesModal.modal('show');
    }
  }
});
