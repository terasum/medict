import { SyncMainAPI } from './service.renderer.manifest';
const ret = SyncMainAPI.syncMessage('myhello');
console.log(`[render-rpc]: syncMessage | ret: ${ret}`);
