import { StubWindow } from './mainsvc/StubWindow';
import { StubMessage } from './mainsvc/StubMessage';
import { StubWordQuery } from './mainsvc/StubWordQuery';
import { StubDictAccessor } from './mainsvc/StubDictAccessor';
import { StubFileOpen } from './mainsvc/StubFileOpen';

const stubWindow = new StubWindow();
const stubMessage = new StubMessage();
const stubWordQuery = new StubWordQuery();
const stubDictAccessor = new StubDictAccessor();
const stubFileOpen = new StubFileOpen();

export const asyncfn = {
  asyncMessage: stubMessage.asyncMessage,
  createSubWindow: stubWindow.createSubWindow,
  entryLinkWord: stubWordQuery.entryLinkWord,
  suggestWord: stubWordQuery.suggestWord,
  findWordPrecisly: stubWordQuery.findWordPrecisly,
  loadDictResource: stubWordQuery.loadDictResource,
};

export const syncfn = {
  syncMessage: stubMessage.syncMessage,
  syncShowOpenDialog: stubFileOpen.syncShowOpenDialog,
  dictAddOne: stubDictAccessor.dictAddOne,
  dictFindOne: stubDictAccessor.dictFindOne,
  dictDeleteOne: stubDictAccessor.dictDeleteOne,
  dictFindAll: stubDictAccessor.dictFindAll,
};
