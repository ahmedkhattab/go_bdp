
var poll = true;
$(document).ready(function(){
  $('.btn-2b').each(function() {
    var input = $(this).children('input');
    if (input.val() === "true")
    $(this).addClass("down");
  });

  $('.btn-2b').click(function(){
    $(this).toggleClass("down");
    var input = $(this).children('input');
    input.val(input.val() === "true" ? "false" : "true");
  });

  $("#mainForm").submit(function(e)
  {
      var postData = $(this).serializeArray();
      var formURL = $(this).attr("action");
      $("#progress_panel").show();
      $.ajax(
      {
          url : formURL,
          type: "POST",
          data : postData,
          success:function(data, textStatus, jqXHR)
          {
              poll = false;
              $('#submit').html('Deploy <span class="glyphicon glyphicon-cloud-upload" aria-hidden="true"></span>').removeClass(
										'loading')
								.addClass(
										'btn-success').data("stop",0);
                    alert("done !");

          },
          error: function(jqXHR, textStatus, errorThrown)
          {
          },
      });
      pollLog();
      e.preventDefault();
      e.stopPropagation(); //STOP default action
  });
$('#submit').click(function(){
  var $this = $(this);
  if($this.data("stop")==1)
    return;
  $this.addClass("loading")
        .removeClass('btn-success')
				.addClass('btn-default')
        .text("deploying...")
        .data("stop",1);

  $("#mainForm").submit(); //Submit  the FORM

});

});
var previous_text = '';
function pollLog() {
  $.ajax(
  {
      url : "/static/log.out",
      type: "GET",
      success:function(data, textStatus, jqXHR)
      {
          var new_text = data.substring(previous_text.length);
          console.log(previous_text.length)
          console.log(new_text)
          if(new_text != '')
            $("#log").append(new_text.split('\n').join('<br/>'));
            $("#log").animate({
              scrollTop: $('#log')[0].scrollHeight}, 1500);
          if(poll === true) {
            setTimeout( function() { previous_text = data; pollLog(); }, 2000);
          }
          else {
            poll = true;
         }
      },
      error: function(jqXHR, textStatus, errorThrown)
      {
      }
  });
}
