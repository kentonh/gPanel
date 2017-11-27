jQuery('._js_back-to-specific-bundle').on('click', function(e){
  e.preventDefault();

  jQuery('.view-log-modal').modal('hide');
  jQuery('.specific-bundle-modal').modal('show');
});

jQuery('.view-log-modal').on('hidden.bs.modal', function(e){
  jQuery('._js_back-to-specific-bundle').removeClass('d-none');
  jQuery('._js_back-to-specific-bundle').addClass('d-none');
});
