$( document ).ready(function() {
  var converter = new showdown.Converter();

  $("p[data-markdown]").each(function(index, el) {
    $el = $(el);
    unsafeHtml = converter.makeHtml($el.data("markdown"));
    saferHtml = filterXSS(unsafeHtml);

    $el
      .removeAttr("data-markdown")
      .html(saferHtml)
      .addClass("markdown")
  });
});
