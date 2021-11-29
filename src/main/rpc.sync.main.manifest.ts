import { getResourceServerPort } from '../main/init.resource.server';
import { createByProc } from '@terasum/electron-call';
import { ClipboardAPI } from './apis/ClipboardAPI';
import { WindowAPI } from './apis/WindowAPI';
import { FileOpenAPI } from './apis/FileOpenAPI';
import * as basicConfigAPI from './apis/BasicConfigAPI';

const stubByMain = createByProc('main', 'error');

const fileOpenApi = new FileOpenAPI();
const clipboardApi = new ClipboardAPI();
const windowApi = new WindowAPI();

// 提供异步API
stubByMain.provide(['renderer', 'worker'], 'WindowApi', windowApi);

// 提供同步API
export const syncfn = {
  syncShowOpenDialog: fileOpenApi.syncShowOpenDialog,
  syncShowOpenDirDialog: fileOpenApi.syncShowOpenDirDialog,
  
  syncGetResourceServerPort: getResourceServerPort,
  syncClipboardWriteText: clipboardApi.syncClipboardWriteText,
  syncGetResourceRootPath: basicConfigAPI.getResourceRootPath,

  syncGetUserResourceRootPath: basicConfigAPI.getUserResourceRootPath,
  syncGetConfigJsonPath: basicConfigAPI.getConfigJsonPath,
  syncGetLoggerFilePath: basicConfigAPI.getLoggerFilePath,
  syncGetWebviewPreliadFilePath: basicConfigAPI.webviewPreloadFilePath,
  syncWritePreloadFile: basicConfigAPI.writePreloadFile,
};
