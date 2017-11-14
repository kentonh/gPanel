jQuery('#loginForm').on('submit', function(e){
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
      jQuery('.index-alert').html('<strong>Login Success: </strong>Redirecting in 2 seconds...');
      jQuery('.index-alert').removeClass('alert-danger').addClass('alert-success').removeClass('d-none');

      setTimeout(function(){
        window.location.href = "/gPanel.html";
      }, 2000);
    }
    else {
      jQuery('.index-alert').html("<strong>Login Error: </strong>" + xhr.response);
      jQuery('.index-alert').removeClass('alert-success').addClass('alert-danger').removeClass('d-none');
    }
  }
});
