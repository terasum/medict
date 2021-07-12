import { FileOpenService } from './svc/FileOpenService';
import { logger } from '../../utils/logger';

const fileOpenService = new FileOpenService();

export class StubFileOpen {
  syncShowOpenDialog(arg: string[]) {
    logger.info('syncShowOpenDialog - arg');
    logger.info(arg);
    return fileOpenService.showOpenDialog(arg);
  }
}
