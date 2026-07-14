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
<!-- out-plugin style and js file -->
<link href="%s.css?dict_id=%s" rel="stylesheet">
<script async src='%s.js?dict_id=%s'></script>

<script lang="javascript">
function __medict_play_sound(mp3url) {
	console.log(mp3url);
	var audioEle = document.createElement("audio");
	audioEle.src = mp3url;
	document.body.appendChild(audioEle);
	audioEle.play();
}

//*************************
// top-inner frame communication
//**************************
var __TOPFRAME_SECURE_ORIGIN__ = "*";
var __MEDICT_DICT_ID__ = "%s";
function __medict_entry_jump(word, dict_id) {
	console.log("[inner frame] jump entry => ", word, dict_id);
	if (window.top){
		window.top.postMessage({"evtype":"__Medict_INNER_FRAME_MSG_EVTP_ENTRY_JUMP", "word":word, "dict_id":dict_id},__TOPFRAME_SECURE_ORIGIN__ )
	}
}

// setup event listener and top frame origin
!(function(){

	let defaultFontSize = 1;
    window.addEventListener('message', function(e) {
        console.log("[inner frame got message] ", e)
		if (e && e.origin && e.origin.startsWith("wails://")){
			if (__TOPFRAME_SECURE_ORIGIN__ !== e.origin){
				console.log("__TOPFRAME_SECURE_ORIGIN__", __TOPFRAME_SECURE_ORIGIN__);
				__TOPFRAME_SECURE_ORIGIN__ = e.origin;
			}
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_ZOOM_OUT"){
			console.log("zoom out event", e);
			defaultFontSize += 0.1;
			document.body.style.fontSize = defaultFontSize + "em";
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_ZOOM_IN"){
			console.log("zoom out event", e);
			defaultFontSize -= 0.1;
			document.body.style.fontSize = defaultFontSize + "em";
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_REFRESH"){
			console.log("refresh event", e);
			window.location.reload();
		}
    })
}())

// Defense-in-depth (#718): intercept any entry:// link the backend replacer did
// NOT rewrite to javascript:__medict_entry_jump(...) (e.g. hrefs that use single
// quotes, or entry IDs the regex still misses). Without this guard the webview
// hands the unknown entry:// scheme to the OS, producing
// "There is no application set to open the URL entry://...".
!(function(){
	var ENTRY_PREFIX = "entry://";
	document.addEventListener("click", function(e) {
		var t = e.target;
		var a = t && t.closest ? t.closest("a") : null;
		if (!a) return;
		var href = a.getAttribute("href") || "";
		if (href.indexOf(ENTRY_PREFIX) !== 0) return;
		e.preventDefault();
		__medict_entry_jump(href.slice(ENTRY_PREFIX.length), __MEDICT_DICT_ID__);
	});
	// Double-click → look up the browser-selected word (#258). The native
	// selection (CJK word-boundary aware) fires on dblclick.
	document.addEventListener("dblclick", function(e) {
		var sel = window.getSelection ? String(window.getSelection()).trim() : "";
		if (sel) {
			window.top.postMessage({"evtype":"__Medict_INNER_FRAME_MSG_EVTP_DBLCLICK_LOOKUP", "word": sel}, __TOPFRAME_SECURE_ORIGIN__);
		}
	});
}());

</script>
</head>
<body>
%s
</body>
</html>
`
