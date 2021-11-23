import { logger } from '../../utils/logger';
import { ClipboardService } from '../mainsvc/ClipboardService';

export class ClipboardApi {
  clipboard: ClipboardService;

  constructor() {
    this.clipboard = new ClipboardService();
  }

  syncClipboardWriteText(arg: { text: string }) {
    logger.info('syncClipboardWriteText - arg');
    logger.info(arg);
    if (!arg || !arg.text) {
      return false;
    }
  }
}
