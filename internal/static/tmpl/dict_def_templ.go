//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
