package tmpl

const WordContainer = `
<html>
<head>
<script lang="javascript">
function setCookie(name, value, days) {
  var expires = "";
  if (days) {
    var date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

  console.log("==== iframe location ====");
  console.log(window.location);
  console.log(window.location.href);
  console.log(document.domain);
  console.log(document.cookie);

  let key1 = "%s";
  let var1 = "%s";
  let key2 = "%s";
  let var2 = "%s";

  setCookie(key1,var1,1);
  console.log("setcookie",key1, var1);


  setCookie(key2,var2,1);
  console.log("setcookie",key2, var2);
  console.log("after setting--------");
  console.log(document.cookie);

</script>

<link rel="stylesheet" type="text/css" href="styles.css"/>
<script type="text/javascript" src="path-to-javascript-file.js"></script>


</head>
<body>
<span>%s:</span><span>%s</span>
</body>
</html>
`
