
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

  var l = Ladda.create($('#submit'));
  $("#mainForm").submit(function(e)
  {
      var postData = $(this).serializeArray();
      var formURL = $(this).attr("action");
      $.ajax(
      {
          url : formURL,
          type: "POST",
          data : postData,
          beforeSend:function(jqXHR) {
            l.start()
          },
          success:function(data, textStatus, jqXHR)
          {
              alert(data)
              l.stop()
          },
          error: function(jqXHR, textStatus, errorThrown)
          {
              //if fails
          }
      });
      e.preventDefault(); //STOP default action
  });
$('#submit').click(function(){
  $("#mainForm").submit(); //Submit  the FORM
});

});
