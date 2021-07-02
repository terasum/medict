import { StubWindow } from './mainsvc/StubWindow';
import { StubMessage } from './mainsvc/StubMessage';
import { StubWordQuery } from './mainsvc/StubWordQuery';
import { StubDictAccessor } from './mainsvc/StubDictAccessor';

const stubWindow = new StubWindow();
const stubMessage = new StubMessage();
const stubWordQuery = new StubWordQuery();
const stubDictAccessor = new StubDictAccessor();

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
  dictAddOne: stubDictAccessor.dictAddOne,
  dictFindOne: stubDictAccessor.dictFindOne,
  dictDeleteOne: stubDictAccessor.dictDeleteOne,
  dictFindAll: stubDictAccessor.dictFindAll,
};
