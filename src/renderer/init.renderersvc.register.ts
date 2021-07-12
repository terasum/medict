import { listeners } from './service.renderer.listener';
import { ipcRenderer } from 'electron';

export function cleanUpListeneres() {
  for (const lis in listeners) {
    if (Object.prototype.hasOwnProperty.call(listeners, lis)) {
      ipcRenderer.removeAllListeners(lis);
      console.log(`🔨 remove renderer async listener: ${lis}`);
    }
  }
}
