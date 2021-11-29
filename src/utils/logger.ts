import { getLoggerFilePath } from '../main/apis/BasicConfigAPI';
import { app } from 'electron';

import log from 'electron-log';
if (app) {
  log.transports.file.resolvePath = () => getLoggerFilePath();
}

export const logger = log.scope('main-process');
