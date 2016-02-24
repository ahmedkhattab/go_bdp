
$(document).ready(function(){
    $('.btn-2b').click(function(){
        $(this).toggleClass("down");
        var input = $(this).children('input');
        input.val(input.val() === "true" ? "false" : "true");
        alert( input.val());
    });
});
