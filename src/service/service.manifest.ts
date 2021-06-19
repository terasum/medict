export const asyncfn = {
  asyncMessage: (event: any, arg: any) => {
    event.sender.send('asynchronous-reply', arg);
  }
}

export const syncfn = {
  syncMessage: (arg: any) => {
    console.log(arg);
    return 'pong';
  }
}