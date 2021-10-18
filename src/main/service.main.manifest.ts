import { StubWindow } from './mainsvc/StubWindow';
import { StubMessage } from './mainsvc/StubMessage';
import { StubWordQuery } from './mainsvc/StubWordQuery';
import { StubDictAccessor } from './mainsvc/StubDictAccessor';
import { StubFileOpen } from './mainsvc/StubFileOpen';
import { StubConfigAccessor } from './mainsvc/StubConfigAccessor';
import { StubTranslate } from './mainsvc/StubTranslate';
import { StubClipboard } from './mainsvc/StubClipboard';

const stubWindow = new StubWindow();
const stubMessage = new StubMessage();
const stubWordQuery = new StubWordQuery();
const stubDictAccessor = new StubDictAccessor();
const stubFileOpen = new StubFileOpen();
const stubConfigAccessor = new StubConfigAccessor();
const stubTranslate = new StubTranslate();
const stubClipboard = new StubClipboard();

export const asyncfn = {
  asyncMessage: stubMessage.asyncMessage,
  createSubWindow: stubWindow.createSubWindow,
  openDevTool: stubWindow.openDevTool,
  openResourceDir: stubWindow.openResourceDir,
  openDictResourceDir: stubWindow.openDictResourceDir,
  openMainProcessLog: stubWindow.openMainProcessLog,
  openUrlOnBrowser: stubWindow.openUrlOnBrowser,
  entryLinkWord: stubWordQuery.entryLinkWord,
  suggestWord: stubWordQuery.suggestWord,
  findWordPrecisly: stubWordQuery.findWordPrecisly,
  loadDictResource: stubWordQuery.loadDictResource,
  asyncBaiduTranslate: stubTranslate.asyncBaiduTranslate,
  asyncGoogleTranslate: stubTranslate.asyncGoogleTranslate,
  asyncYoudaoTranslate: stubTranslate.asyncYoudaoTranslate,
};

export const syncfn = {
  syncMessage: stubMessage.syncMessage,
  syncShowOpenDialog: stubFileOpen.syncShowOpenDialog,
  syncShowMainLoggerPath: stubFileOpen.syncShowMainLoggerPath,
  syncGetResourceRootPath: stubFileOpen.syncGetResourceRootPath,
  syncGetWebviewPreliadFilePath: stubFileOpen.syncGetWebviewPreliadFilePath,
  syncShowComfirmMessageBox: stubWindow.syncShowComfirmMessageBox,
  dictAddOne: stubDictAccessor.dictAddOne,
  dictFindOne: stubDictAccessor.dictFindOne,
  dictDeleteOne: stubDictAccessor.dictDeleteOne,
  dictFindAll: stubDictAccessor.dictFindAll,
  loadTranslateApiConfig: stubConfigAccessor.loadTranslateApiConfig,
  saveTranslateBaiduApiConfig: stubConfigAccessor.saveTranslateBaiduApiConfig,
  saveTranslateYoudaoApiConfig: stubConfigAccessor.saveTranslateYoudaoApiConfig,
  clipboardWriteText: stubClipboard.syncClipboardWriteText,
};
