let amount = 3;

$(document).on('keydown', '.answer', function(event) {
    let id = this.id;
    let key = event.which;

    if(key !== 8 && this.value.length === 0 && id.endsWith(amount)) {
        amount++;
        $('#answers-wrapper').append("<div class=\"input-field\" id=\"wrapper-answer" + amount + "\">\n" +
                                        "<input id=\"answer" + amount + "\" type=\"text\" class=\"validate answer\" name=\"answer" + amount + "\">\n" +
                                        "<label for=\"answer" + amount + "\">Poll option</label>\n" +
                                     "</div>")

        $('#answers-amount').val(amount)
    }
});
