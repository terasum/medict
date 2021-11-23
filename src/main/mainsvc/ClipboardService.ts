import { logger } from '../../utils/logger';
import { clipboard } from 'electron';

export class ClipboardService {
  clipboardWriteText(text: string) {
    if (!text || text.length == 0) {
        return false;
    }
    clipboard.writeText(text, 'selection')
    return true;
  }
}
