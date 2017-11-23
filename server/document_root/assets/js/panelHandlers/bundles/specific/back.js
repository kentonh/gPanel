jQuery('._js_back-to-bundle-menu').on('click', function(e){
  e.preventDefault();

  jQuery('.specific-bundle-modal').modal('hide');
  jQuery('.manage-bundles-modal').modal('show');
});
