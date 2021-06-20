

export const asyncfn = {
  asyncMessage: (event: any, arg: any) => {
    event.sender.send('asynchronous-reply', arg);
  },
  createSubWindow: (event: any, arg: any) => {
    event.sender.send('createSubWindow', arg);
  }
}

export const syncfn = {
  syncMessage: (arg: any) => {
    console.log(arg);
    return 'pong';
  }
}