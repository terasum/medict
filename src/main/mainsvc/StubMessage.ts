export class StubMessage {
  asyncMessage(event: any, arg: any) {
    event.sender.send('asynchronous-reply', arg);
  }
  syncMessage(arg: any) {
    console.log(arg);
    return 'pong';
  }
}
