import { FileOpenService } from './svc/FileOpenService';
import { logger } from '../../utils/logger';
import { clipboard } from 'electron';

export class StubClipboard {
  syncClipboardWriteText(arg: {text: string}) {
    logger.info('syncClipboardWriteText - arg');
    logger.info(arg);
    if (arg && arg.text && arg.text.length > 0) {
        clipboard.writeText(arg.text, 'selection')
        return true;
    }
    return false;
  }
}
