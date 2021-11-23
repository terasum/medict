import { FileOpenService } from '../mainsvc/FileOpenService';
import { logger } from '../../utils/logger';
import { getUserResourceRootPath, getResourceRootPath, getLoggerFilePath, webviewPreloadFilePath } from '../../config/config';


export class FileOpenApi {
  fileOpenService: FileOpenService;

  constructor() {
    this.fileOpenService = new FileOpenService();
  }
  syncShowOpenDialog(arg: {fileExtensions: string[] | undefined, multiFile: boolean}) {
    logger.info('syncShowOpenDialog - arg');
    logger.info(arg);
    return this.fileOpenService.showOpenDialog(arg);
  }
  syncShowMainLoggerPath(arg?: any) {
    return getLoggerFilePath();
  }

  syncGetResourceRootPath(arg?: any) {
    return getResourceRootPath();
  }
  
  syncGetUserResourceRootPath(arg?: any) {
    return getUserResourceRootPath();
  }

  syncGetWebviewPreliadFilePath(arg?:any) {
    return webviewPreloadFilePath();
  }
}

