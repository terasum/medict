import { dialog } from 'electron';
import { logger } from '../../../utils/logger';

export class FileOpenService {
  showOpenDialog(fileExtensions: string[] | undefined) {
    logger.info({ fileExtensions });
    if (!fileExtensions) {
      fileExtensions = ['mdd', 'mdx'];
    }
    return dialog.showOpenDialogSync({
      properties: ['openFile', 'noResolveAliases', 'dontAddToRecent'],
      message: '选择文件',
      filters: [{ name: 'Custom File Type', extensions: fileExtensions }],
    });
  }
}
