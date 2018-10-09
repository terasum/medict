class CommuniMsg {
  constructor (msgType, data, code) {
    this.msgType = msgType
    this.data = data
    this.code = code || 0
  }
}

export default CommuniMsg
