import { shell, dialog, MessageBoxSyncOptions } from 'electron'; // deconstructing assignment

import { getResourceRootPath, getLoggerFilePath } from '../../config/config';
import path from 'path';
import fs from 'fs';

/**
 * WindowService 创建新窗口服务
 */
export class StubWindow {
  /**
   * createSubWindow 创建新的子窗口
   * @param event 事件源
   * @param arg 窗口选项
   */
  createSubWindow(
    event: any,
    arg: {
      height: number;
      width: number;
      titleBarStyle: string;
      nodeIntegration: boolean;
      contextIsolation: boolean;
      html: string;
    }
  ) {
    event.sender.send('createSubWindow', arg);
  }
  openDevTool(event: any) {
    event.sender.send('openDevTool');
  }
  openResourceDir(event: any) {
    console.log(`openResourceDir ${getResourceRootPath()}`);
    shell.openPath(getResourceRootPath()); // Open the given file in the desktop's default manner.
  }
  openDictResourceDir(event: any, dictid: string) {
    const fpath = path.resolve(getResourceRootPath(), dictid);
    if (fs.existsSync(fpath)) {
      shell.openPath(fpath); // Open the given file in the desktop's default manner.
    }
  }
  openMainProcessLog(event: any) {
    console.log(`openMainProcessLog ${getLoggerFilePath()}`);
    const logpath = getLoggerFilePath();
    if (fs.existsSync(logpath)) {
      shell.openPath(logpath); // Open the given file in the desktop's default manner.
    }
  }
  openUrlOnBrowser(event: any, url: string) {
    console.log(`openurl ${url}`);
    if (url && url.length > 0 && url.startsWith('https://')) {
      shell.openExternal(url);
    }
  }
  syncShowComfirmMessageBox(args: MessageBoxSyncOptions) {
    console.log('showComfirmMessageBox args: ');
    console.log(args);
    return dialog.showMessageBoxSync(args);
  }
}
