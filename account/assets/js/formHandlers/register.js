jQuery('#registerForm').on('submit', function(e){
  e.preventDefault();

  var formData = {};
  for(var y = 0, yy = this.length; y < yy; y++) {
    var input = this[y];
    if(input.name) {
      formData[input.name] = input.value;
    }
  }

  var xhr = new XMLHttpRequest();
  xhr.open(jQuery(this).attr('method'), jQuery(this).attr('action'), true);
  xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
  xhr.send(JSON.stringify(formData));

  xhr.onloadend = function() {
    if(xhr.status == 200 || xhr.status == 204) {
      jQuery('.index-alert').html('<strong>Register Success: </strong>You may now login.');
      jQuery('.index-alert').removeClass('alert-danger').addClass('alert-success').removeClass('d-none');
    }
    else {
      jQuery('.index-alert').html("<strong>Register Error: </strong>" + xhr.response);
      jQuery('.index-alert').removeClass('alert-success').addClass('alert-danger').removeClass('d-none');
    }
  }
});
