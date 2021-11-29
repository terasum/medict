export const preloadContent = `
const { ipcRenderer } = require('electron');
console.warn('=== preload electron [sandbox] ===');

window.addEventListener('message', function (event) {
   console.log('---- preload listenning message -----');
   console.log(event);
   if (!event.data || !event.data.channel || !event.data.payload) {
     return;
   }
   console.log(
     'send to renderer-process [\${event.data.channel}|\${event.data.payload}]'
   );
   ipcRenderer.sendToHost(event.data.channel, event.data.payload);
 });
`;
