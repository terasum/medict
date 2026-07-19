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

	let contentZoom = 100;
	function __medict_apply_content_zoom(value) {
		contentZoom = Math.min(160, Math.max(80, Number(value) || 100));
		document.documentElement.style.zoom = String(contentZoom / 100);
	}
    window.addEventListener('message', function(e) {
        console.log("[inner frame got message] ", e)
		if (e && e.origin && e.origin.startsWith("wails://")){
			if (__TOPFRAME_SECURE_ORIGIN__ !== e.origin){
				console.log("__TOPFRAME_SECURE_ORIGIN__", __TOPFRAME_SECURE_ORIGIN__);
				__TOPFRAME_SECURE_ORIGIN__ = e.origin;
			}
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_ZOOM_OUT"){
			__medict_apply_content_zoom(contentZoom - 10);
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_ZOOM_IN"){
			__medict_apply_content_zoom(contentZoom + 10);
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_SET_ZOOM"){
			__medict_apply_content_zoom(e.data.scale);
		}
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_REFRESH"){
			console.log("refresh event", e);
			window.location.reload();
		}
		// #783: apply user CSS override (live preview from the CSS editor)
		if (e && e.data && e.data.evtype === "__Medict_TOP_WIN_MSG_EVTP_APPLY_USER_CSS"){
			var ex = document.getElementById("medict-user-css");
			if (ex) { ex.remove(); }
			if (e.data.css) {
				var st = document.createElement("style");
				st.id = "medict-user-css";
				st.textContent = e.data.css;
				document.head.appendChild(st);
			}
		}
    })
}())

// Defense-in-depth (#718): intercept entry:// links and provide macOS Dictionary-
// style word interaction: hover underlines a segmented word, click queries it.
!(function(){
	var ENTRY_PREFIX = "entry://";
	var __wordHint = null;
	var __hoveredWord = "";
	var __hoverFrame = 0;
	var __lastPointer = null;
	var __hoveredNode = null;
	var __hoveredStart = -1;
	var __hoveredEnd = -1;
	var __originalCursor = document.documentElement.style.cursor;
	var __wordSegmenter = typeof Intl !== "undefined" && Intl.Segmenter
		? new Intl.Segmenter(undefined, { granularity: "word" })
		: null;

	function __medict_clearWordHint() {
		if (__wordHint) {
			__wordHint.remove();
			__wordHint = null;
		}
		__hoveredWord = "";
		__hoveredNode = null;
		__hoveredStart = -1;
		__hoveredEnd = -1;
		document.documentElement.style.cursor = __originalCursor;
	}

	function __medict_isCurrentWord(hit) {
		return hit && hit.word === __hoveredWord && hit.node === __hoveredNode
			&& hit.start === __hoveredStart && hit.end === __hoveredEnd;
	}

	function __medict_caretRangeFromPoint(x, y) {
		if (document.caretRangeFromPoint) {
			return document.caretRangeFromPoint(x, y);
		}
		if (document.caretPositionFromPoint) {
			var pos = document.caretPositionFromPoint(x, y);
			if (!pos) return null;
			var range = document.createRange();
			range.setStart(pos.offsetNode, pos.offset);
			range.collapse(true);
			return range;
		}
		return null;
	}

	function __medict_wordRangeAtPoint(x, y) {
		var caret = __medict_caretRangeFromPoint(x, y);
		if (!caret || !caret.startContainer || caret.startContainer.nodeType !== 3) return null;
		var node = caret.startContainer;
		var parent = node.parentElement;
		if (!parent || parent.closest("a,button,input,textarea,select,option,script,style")) return null;
		var text = node.textContent || "";
		if (!text) return null;
		var offset = Math.min(caret.startOffset, text.length - 1);
		if (offset < 0 || /\s/.test(text.charAt(offset))) return null;
		var start = -1, end = -1, word = "";

		if (__wordSegmenter) {
			var iterator = __wordSegmenter.segment(text)[Symbol.iterator]();
			var step = iterator.next();
			while (!step.done) {
				var part = step.value;
				if (part.isWordLike && offset >= part.index && offset < part.index + part.segment.length) {
					start = part.index;
					end = part.index + part.segment.length;
					word = part.segment;
					break;
				}
				step = iterator.next();
			}
		}

		if (start < 0) {
			var latin = /[A-Za-z0-9'-]/;
			var ch = text.charAt(offset);
			if (latin.test(ch)) {
				start = offset;
				end = offset + 1;
				while (start > 0 && latin.test(text.charAt(start - 1))) start--;
				while (end < text.length && latin.test(text.charAt(end))) end++;
				word = text.slice(start, end);
			} else if (/[\u3400-\uFAFF]/.test(ch)) {
				start = offset;
				end = offset + 1;
				word = ch;
			}
		}

		word = word.trim();
		if (!word || start < 0 || end <= start) return null;
		var range = document.createRange();
		range.setStart(node, start);
		range.setEnd(node, end);
		var rects = Array.from(range.getClientRects()).filter(function(rect) {
			return rect.width > 0 && rect.height > 0 && x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom;
		});
		if (rects.length === 0) return null;
		return { word: word, range: range, node: node, start: start, end: end };
	}

	function __medict_showWordHint(hit) {
		__medict_clearWordHint();
		var hint = document.createElement("div");
		hint.setAttribute("aria-hidden", "true");
		hint.style.cssText = "position:fixed;inset:0;pointer-events:none;z-index:2147483647";
		Array.from(hit.range.getClientRects()).forEach(function(rect) {
			var line = document.createElement("span");
			line.style.cssText = "position:fixed;pointer-events:none;border-bottom:1px solid currentColor;left:" + rect.left + "px;top:" + (rect.bottom - 1) + "px;width:" + rect.width + "px;height:1px";
			hint.appendChild(line);
		});
		document.body.appendChild(hint);
		__wordHint = hint;
		__hoveredWord = hit.word;
		__hoveredNode = hit.node;
		__hoveredStart = hit.start;
		__hoveredEnd = hit.end;
		document.documentElement.style.cursor = "pointer";
	}

	document.addEventListener("click", function(e) {
		var t = e.target;
		var a = t && t.closest ? t.closest("a") : null;
		if (a) {
			var href = a.getAttribute("href") || "";
			if (href.indexOf(ENTRY_PREFIX) === 0) {
				e.preventDefault();
				__medict_entry_jump(href.slice(ENTRY_PREFIX.length), __MEDICT_DICT_ID__);
			}
			return;
		}
		var hit = __medict_wordRangeAtPoint(e.clientX, e.clientY);
		if (!__medict_isCurrentWord(hit)) return;
		window.top.postMessage({"evtype":"__Medict_INNER_FRAME_MSG_EVTP_CLICK_LOOKUP", "word": hit.word}, __TOPFRAME_SECURE_ORIGIN__);
		__medict_clearWordHint();
	});

	document.addEventListener("mousemove", function(ev) {
		__lastPointer = { x: ev.clientX, y: ev.clientY };
		if (__hoverFrame) return;
		__hoverFrame = requestAnimationFrame(function() {
			__hoverFrame = 0;
			if (!__lastPointer) return;
			var hit = __medict_wordRangeAtPoint(__lastPointer.x, __lastPointer.y);
			if (!hit) {
				__medict_clearWordHint();
				return;
			}
			if (!__medict_isCurrentWord(hit)) __medict_showWordHint(hit);
		});
	});
	document.addEventListener("mouseleave", __medict_clearWordHint);
	window.addEventListener("scroll", __medict_clearWordHint, true);
	window.addEventListener("keydown", function(ev){ if (ev.key === "Escape") __medict_clearWordHint(); });
}());

</script>
</head>
<body>
%s
</body>
</html>
`
