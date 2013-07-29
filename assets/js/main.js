

function success (data, status, xhr) {
  alert(status)
}
function post() {
  var text = $("#text")[0].value
  var obj = {}
  obj.text = text
  $.post("/scard/something", JSON.stringify(obj), "application/json")
}
