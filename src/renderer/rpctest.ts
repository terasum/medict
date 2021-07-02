import { MainProcSyncAPI } from './service.renderer.manifest';
const ret = MainProcSyncAPI.syncMessage('myhello');
console.log(`[render-rpc]: syncMessage | ret: ${ret}`);
