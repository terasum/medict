package tmpl

const WordDefinitionTempl = `
<html>
<head>
<style>
body{
	padding: 5px 6px 5px 6px;
}
</style>
<script lang="javascript">
function __medict_play_sound(mp3url) {
	console.log(mp3url);
	var audioEle = document.createElement("audio");
	audioEle.src = mp3url;
	document.body.appendChild(audioEle);
	audioEle.play();
}
</script>
</head>
<body>
%s
</body>
</html>
`
