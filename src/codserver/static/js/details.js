var poll = true;

function refreshToggleBtns() {
  $('.btn-cstm').each(function() {
    var input = $(this).children('input');
    var tick = $(this).children('.checkmark');
    var stopped = $(this).children('.stopped');
    if (input.val() === "true") {
      $(this).addClass("running");
      tick.show();
      stopped.hide();
    }
  });
}

function validate() {
  var c = 0;
  $('.btn-2b').each(function() {
    var input = $(this).children('input');
    if (input.val() === "true") {
      c++;
    }
  });
  if (c === 0)
    return "No componets were chosen for deployment";
  else
    return "";
}

function resetSubmitBtn() {
  $('#submit').html(
      'Deploy <span class="glyphicon glyphicon-cloud-upload" aria-hidden="true"></span>'
    ).removeClass(
      'loading')
    .addClass(
      'btn-success').data("stop", 0);
}

$(document).ready(function() {
  refreshToggleBtns();
  $('[data-toggle="popover"]').popover()

  $('.btn-2b').click(function() {
    $(this).toggleClass("down");
    var input = $(this).children('input');
    input.val(input.val() === "true" ? "false" : "true");
  });

  $("#mainForm").submit(function(e) {

    var postData = $(this).serializeArray();
    var formURL = $(this).attr("action");
    setInterval(type, 600);
    $('#progressModal').modal('show');

    pollLog();
  });

  $('#submit').click(function(e) {
    var err = validate();
    if (err != "") {
      alert(err);
      e.preventDefault();
      e.stopPropagation();
      resetSubmitBtn();
      return;
    }
    var $this = $(this);
    if ($this.data("stop") == 1)
      return;
    $this.addClass("loading")
      .removeClass('btn-success')
      .addClass('btn-default')
      .text("deploying...")
      .data("stop", 1);
  });

});
var previous_text = '';

function pollLog() {
  $.ajax({
    url: "/static/log.out",
    type: "GET",
    success: function(data, textStatus, jqXHR) {
      var new_text = data.substring(previous_text.length);
      console.log(previous_text.length)
      console.log(new_text)
      if (new_text != '')
        $("#log").append(new_text.split('\n').join('<br/>'));
      $("#log").animate({
        scrollTop: $('#log')[0].scrollHeight
      }, 1500);
      if (poll === true) {
        setTimeout(function() {
          previous_text = data;
          pollLog();
        }, 1000);
      } else {
        setTimeout(function() {
          $('#progressModal').modal('toggle');
        }, 2000);
        poll = true;
      }
    },
    error: function(jqXHR, textStatus, errorThrown) {}
  });
}

var dots = 0;

function type() {
  if (dots < 3) {
    $('#dots').append('.');
    dots++;
  } else {
    $('#dots').html('');
    dots = 0;
  }
}
