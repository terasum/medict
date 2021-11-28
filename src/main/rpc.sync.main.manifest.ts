import { FileOpenApi } from './apis/FileOpenApi';
import { getResourceServerPort } from '../main/init.resource.server';
import { createByProc } from '@terasum/electron-call';
import { ClipboardApi } from './apis/ClipboardApi';
import { WindowApi } from './apis/WindowApi';

const stubByMain = createByProc('main', 'error');

const fileOpenApi = new FileOpenApi();
const clipboardApi = new ClipboardApi();
const windowApi = new WindowApi();

// 提供异步API
stubByMain.provide(['renderer', 'worker'], 'WindowApi', windowApi);

// 提供同步API
export const syncfn = {
  syncShowOpenDialog: fileOpenApi.syncShowOpenDialog,
  syncShowMainLoggerPath: fileOpenApi.syncShowMainLoggerPath,
  syncGetResourceRootPath: fileOpenApi.syncGetResourceRootPath,
  syncGetWebviewPreliadFilePath: fileOpenApi.syncGetWebviewPreliadFilePath,
  syncGetResourceServerPort: getResourceServerPort,
  syncClipboardWriteText: clipboardApi.syncClipboardWriteText,
};
