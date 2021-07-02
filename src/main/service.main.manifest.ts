import _ from 'lodash';
import { WindowService } from './mainsvc/WindowService';
import { MessageService } from './mainsvc/MessageService';
import { WordService } from './mainsvc/WordService';

const windowService = new WindowService();
const messageService = new MessageService();
const wordService = new WordService();

export const asyncfn = {
  asyncMessage: messageService.asyncMessage,
  createSubWindow: windowService.createSubWindow,
  entryLinkWord: wordService.entryLinkWord,
  suggestWord: wordService.suggestWord,
  findWordPrecisly: wordService.findWordPrecisly,
  loadDictResource: wordService.loadDictResource,
};

export const syncfn = {
  syncMessage: messageService.syncMessage,
};
