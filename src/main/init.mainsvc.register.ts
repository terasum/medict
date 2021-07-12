import { syncfn, asyncfn } from './service.main.manifest';
import { ipcMain } from 'electron';
import { logger } from '../utils/logger';

export function registerServices() {
  logger.info(
    '\n================  REGISTER MAIN SERVICE ===================\n'
  );
  for (const fnName in syncfn) {
    if (Object.prototype.hasOwnProperty.call(syncfn, fnName)) {
      const fn = syncfn[fnName];
      logger.info('🔧 register main process sync  service: %s', fnName);
      ipcMain.on(fnName, function(event, args) {
        logger.info(
          `[main-rpc:sync]: ================= [${fnName}] =============== START`
        );
        logger.info(args, `[main-rpc:sync]: ${fnName}| arg:`);
        const ret = fn(args);
        logger.info(ret, `[main-rpc:sync]: ${fnName}| ret:`);
        event.returnValue = ret;
        logger.info(
          `[main-rpc:sync]: ================= [${fnName}] =============== END`
        );
      });
    }
  }
  for (const fnName in asyncfn) {
    if (Object.prototype.hasOwnProperty.call(asyncfn, fnName)) {
      const fn = asyncfn[fnName];
      logger.info('🔧 register main process async service: %s', fnName);
      ipcMain.on(fnName, function(event, args) {
        logger.info(
          `[main-rpc:async]: ================= [${fnName}] =============== START`
        );
        logger.info(args, `[main-rpc:asyn](start): ${fnName} - args:`);
        fn(event, args);
        logger.info(
          `[main-rpc:async]: ================= [${fnName}] =============== END`
        );
      });
    }
  }
}
