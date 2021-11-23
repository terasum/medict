import { dialog } from 'electron';
import { logger } from '../../utils/logger';

export class FileOpenService {
  showOpenDialog(arg: {
    fileExtensions: string[] | undefined;
    multiFile: boolean;
  }) {
    logger.info(arg.fileExtensions);
    if (!arg.fileExtensions) {
      arg.fileExtensions = ['mdd', 'mdx'];
    }
    let prop = ['openFile', 'noResolveAliases', 'dontAddToRecent'] as (
      | 'openFile'
      | 'noResolveAliases'
      | 'dontAddToRecent'
      | 'multiSelections'
      | 'openDirectory'
      | 'showHiddenFiles'
      | 'createDirectory'
      | 'promptToCreate'
      | 'treatPackageAsDirectory'
    )[];
    if (arg.multiFile) {
      prop = [
        'openFile',
        'noResolveAliases',
        'dontAddToRecent',
        'multiSelections',
      ];
    }
    return dialog.showOpenDialogSync({
      properties: prop,
      message: '选择文件',
      filters: [{ name: 'Custom File Type', extensions: arg.fileExtensions }],
    });
  }
}
