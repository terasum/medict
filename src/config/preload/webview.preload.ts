export const preloadContent = `
const { ipcRenderer } = require('electron');

console.warn('=== preload electron [sandbox] ===');

// 监听 main-process 发回来的 结果，格式是 {keyText:"", definition:""}
ipcRenderer.on('onFindWordPrecisly', (event, args) => {
  console.log('------ webview listener[onFindWordPrecisly] -----');
  console.log(args);
  return ipcRenderer.sendToHost('onFindWordPrecisly', args);
});

// 主要处理点击 entry://之后的逻辑
// 将会把需要查询的词发送到 main-process
window.addEventListener('message', function (event) {
  console.log('---- preload listenning message -----');
  console.log(event.data);
  if (event.data && event.data.channel && event.data.payload) {
    console.log(
      'send to main-process [\${event.data.channel}|\${event.data.payload}]'
    );
    ipcRenderer.send(event.data.channel, event.data.payload);
  }
});
`;
