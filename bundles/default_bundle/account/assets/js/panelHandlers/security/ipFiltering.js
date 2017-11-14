var ipModal = jQuery('.ip-filter-modal');

jQuery('._js_ip-filtering-open').on('click', function(e){
  e.preventDefault();

  var title;
  switch(jQuery(this).attr('data')) {
    case "general":
      title = "General";
      break;
    case "maintenance":
      title = "Maintenance Mode";
      break;
    default:
      alert("Error, refresh and try again. If problem persists contact server administrator.");
      return;
  }
  title += " IP Filtering";

  ipModal.find('.modal-title').html(title);
  ipModal.modal('show');
});
