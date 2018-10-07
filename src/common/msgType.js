/*********************************
 *      MAIN <=> RENDERER        *
 *********************************/
const MsgToMain = 'MessageToMain'
const MsgToBackground = 'MessageToBackground'

/*********************************
 *     MAINFRAME <=> BACKGROUND  *
 *********************************/

const SubMsgQueryBackground = 'SubMessageQueryWordBackground'
const SubMsgQueryResponse = 'SubMessageQueryWordResponse'
const SubMsgQueryListResponse = 'SubMessageQueryWordListResponse'
const SubMsgPing = 'SubMessagePing'
const SubMsgPong = 'SubMessagePong'

/*********************************
 *     MAINFRAME <=> BACKGROUND  *
 *********************************/

const BGWorkerMsgToWorker = 'bgWorkerMsgToWorker'
const BGWorkerMsgToMain = 'bgWorkerMsgToMain'

const BGWorkerSubMsgQuery = 'BGWorkerSubMsgQuery'
const BGWorkerSubMsgResponse = 'BGWorkerSubMsgResponse'

export default {
  MsgToMain,
  MsgToBackground,
  SubMsgQueryBackground,
  SubMsgQueryResponse,
  SubMsgQueryListResponse,
  SubMsgPing,
  SubMsgPong,

  BGWorkerMsgToWorker,
  BGWorkerMsgToMain,
  BGWorkerSubMsgQuery,
  BGWorkerSubMsgResponse
}
