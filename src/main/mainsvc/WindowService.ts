/**
 * WindowService 创建新窗口服务
 */
export class WindowService {
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
}
