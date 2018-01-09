var sshModal = jQuery(".manage-ssh-modal");

jQuery('._js_manage-ssh').on('click', function(e){
  e.preventDefault();
  sshModal.modal('show');
});
