jQuery('#logoutForm').on('submit', function(e){
  e.preventDefault();

  var check = confirm('Are you sure you want to logut?');

  if(check) {
    var xhr = new XMLHttpRequest();
    xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
    xhr.send();

    xhr.onloadend = function() {
      if(xhr.status == 200 || xhr.status == 204) {
        window.location.href = '/';
      }
      else {
        alert('An error has occurred. Please contact your server\'s administrator.');
      }
    }
  }
});
