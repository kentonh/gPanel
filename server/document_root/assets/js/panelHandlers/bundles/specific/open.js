jQuery(document).on('click', '._js_specific-bundle', function(e){
  e.preventDefault();

  var name = jQuery(this).attr('data');

  jQuery('.manage-bundles-modal').modal('hide');
  jQuery('.specific-bundle-modal').find('.modal-title').html("Bundle \"" + name + "\" Management");
  jQuery('.specific-bundle-modal').modal('show');
});
