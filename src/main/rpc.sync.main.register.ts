import { syncfn, asyncfn } from './rpc.sync.main.manifest';
import { ipcMain } from 'electron';
import { logger } from '../utils/logger';

export function registerServices() {
  logger.info(
    '\n================  REGISTER MAIN SERVICE ===================\n'
  );
  for (const fnName in syncfn) {
    if (Object.prototype.hasOwnProperty.call(syncfn, fnName)) {
      const fn = syncfn[fnName];
      logger.debug('🔧 register main process sync  service: %s', fnName);
      ipcMain.on(fnName, function(event, args) {
        logger.debug(
          `[main-rpc:sync]: ================= [${fnName}] =============== START`
        );
        logger.debug(`[main-rpc:sync]: ${fnName}| arg:`, args);
        const ret = fn(args);
        logger.debug(`[main-rpc:sync]: ${fnName}| ret:`, ret);
        event.returnValue = ret;
        logger.debug(
          `[main-rpc:sync]: ================= [${fnName}] =============== END`
        );
      });
    }
  }
  for (const fnName in asyncfn) {
    if (Object.prototype.hasOwnProperty.call(asyncfn, fnName)) {
      const fn = asyncfn[fnName];
      logger.debug('🔧 register main process async service: %s', fnName);
      ipcMain.on(fnName, function(event, args) {
        logger.debug(
          `[main-rpc:async]: ================= [${fnName}] =============== START`
        );
        logger.debug(`[main-rpc:asyn](start): ${fnName} - args:`);
        logger.debug(args);
        fn(event, args);
        logger.debug(
          `[main-rpc:async]: ================= [${fnName}] =============== END`
        );
      });
    }
  }
}