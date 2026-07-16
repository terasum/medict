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
// hover 查词配置:由父窗经 SETUP 消息下发(#780)。默认开、400ms。
var __medictHoverEnabled = true;
var __medictHoverDelayMs = 400;
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
		// SETUP 携带 hover 设置(父侧经 #777 preferences 控制,免每请求注入模板)
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG__EVTY_SETUP__"){
			if (typeof e.data.hoverEnabled === "boolean") { __medictHoverEnabled = e.data.hoverEnabled; }
			if (typeof e.data.hoverDelayMs === "number" && e.data.hoverDelayMs > 0) { __medictHoverDelayMs = e.data.hoverDelayMs; }
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
	// ===== 悬停弹窗取词(#780)=====
	// 取光标位置的「词」:caretRangeFromPoint 得文本节点 + offset;CJK(0x3400-0xFAFF)
	// 按单字返回,拉丁/数字/连字符按词边界扩展。
	function __medict_wordAtPoint(x, y) {
		var rng = document.caretRangeFromPoint ? document.caretRangeFromPoint(x, y) : null;
		if (!rng || !rng.startContainer) { return ""; }
		var node = rng.startContainer;
		if (node.nodeType !== 3 /* TEXT_NODE */) { return ""; }
		var text = node.textContent || "";
		var off = rng.startOffset;
		if (off < 0 || off >= text.length) { return ""; }
		var ch = text.charCodeAt(off);
		if (ch >= 0x3400 && ch <= 0xFAFF) { return text.charAt(off); } // CJK 单字
		var re = /[A-Za-z0-9'-]/;
		if (!re.test(text.charAt(off))) { return ""; }
		var s = off, end = off;
		while (s > 0 && re.test(text.charAt(s - 1))) { s--; }
		while (end < text.length - 1 && re.test(text.charAt(end + 1))) { end++; }
		return text.slice(s, end + 1);
	}
	var __hoverTimer = null;
	var __lastHoverWord = "";
	function __medict_hoverLeave() {
		if (__hoverTimer) { clearTimeout(__hoverTimer); __hoverTimer = null; }
		if (__lastHoverWord !== "") {
			__lastHoverWord = "";
			window.top.postMessage({"evtype":"__Medict_INNER_FRAME_MSG_EVTP_HOVER_LEAVE"}, __TOPFRAME_SECURE_ORIGIN__);
		}
	}
	document.addEventListener("mousemove", function(ev) {
		if (!__medictHoverEnabled) { return; }
		var x = ev.clientX, y = ev.clientY;
		if (__hoverTimer) { clearTimeout(__hoverTimer); }
		__hoverTimer = setTimeout(function() {
			__hoverTimer = null;
			var w = __medict_wordAtPoint(x, y);
			if (w && w !== __lastHoverWord) {
				__lastHoverWord = w;
				window.top.postMessage({"evtype":"__Medict_INNER_FRAME_MSG_EVTP_HOVER_LOOKUP", "word": w, "clientX": x, "clientY": y}, __TOPFRAME_SECURE_ORIGIN__);
			}
		}, __medictHoverDelayMs);
	});
	document.addEventListener("mouseleave", __medict_hoverLeave);
	window.addEventListener("scroll", __medict_hoverLeave, true);
	window.addEventListener("keydown", function(ev){ if (ev.key === "Escape") __medict_hoverLeave(); });
}());

</script>
</head>
<body>
%s
</body>
</html>
`
