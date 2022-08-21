console.log('👋 This message is being logged by "renderer.init.ts", included via webpack');


(function errorListen() {
  window.onerror = function (error, url, line) {
    // ipcRenderer.send('errorInWindow', { error, url, line });
  };
})();

