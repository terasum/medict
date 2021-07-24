import { FileOpenService } from './svc/FileOpenService';
import { logger } from '../../utils/logger';
import { getResourceRootPath, getLoggerFilePath } from '../../config/config';

const fileOpenService = new FileOpenService();

export class StubFileOpen {
  syncShowOpenDialog(arg: {fileExtensions: string[] | undefined, multiFile: boolean}) {
    logger.info('syncShowOpenDialog - arg');
    logger.info(arg);
    return fileOpenService.showOpenDialog(arg);
  }
  syncShowMainLoggerPath(arg?: any) {
    return getLoggerFilePath();
  }
  syncGetResourceRootPath(arg?: any) {
    return getResourceRootPath();
  }
}
