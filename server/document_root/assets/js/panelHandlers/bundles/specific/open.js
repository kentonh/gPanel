var specificBundleModal = jQuery('.specific-bundle-modal');

jQuery(document).on('click', '._js_specific-bundle', function(e){
  e.preventDefault();

  var name = jQuery(this).attr('data');

  jQuery('.manage-bundles-modal').modal('hide');
  specificBundleModal.find('.modal-title').html("Bundle \"" + name + "\" Management");
  specificBundleModal.attr('data', name);

  getPublicServerStatus(name);

  specificBundleModal.modal('show');
});
