import { syncfn, asyncfn } from './service.main.manifest';
import { ipcMain } from 'electron';

export function registerServices() {
  console.log(
    '\n================  REGISTER MAIN SERVICE ===================\n'
  );
  for (const fnName in syncfn) {
    if (Object.prototype.hasOwnProperty.call(syncfn, fnName)) {
      const fn = syncfn[fnName];
      console.log('🔧 register main process sync  service:', fnName);
      ipcMain.on(fnName, function(event, args) {
        console.log(`[main-rpc:sync]: ${fnName}| arg: ${args}`);
        const ret = fn(args);
        console.log(`[main-rpc:sync]: ${fnName}| ret: ${ret}`);
        event.returnValue = ret;
      });
    }
  }
  for (const fnName in asyncfn) {
    if (Object.prototype.hasOwnProperty.call(asyncfn, fnName)) {
      const fn = asyncfn[fnName];
      console.log('🔧 register main process async service:', fnName);
      ipcMain.on(fnName, function(event, args) {
        console.log(`[main-rpc:asyn](start): ${fnName} - args:`);
        console.log(args);
        fn(event, args);
        console.log(`[main-rpc:asyn](end): ${fnName}`);
      });
    }
  }
}
