const ipc = require('electron').ipcRenderer;

interface msgPayload {
    data: any
    channel: string
    datetime: Date
    origin: string
}

const ChannelMain = 'ipc:main';
const ChannelMainWeb = 'ipc:mainweb';
const ChannelWorker = 'ipc:worker';


// 发送消息至主进程
export function sendMain(channel: string, arg: any) {
    ipc.send(ChannelMain,
        {
            channel: channel,
            data: arg,
            datetime: new Date(),
            origin: 'worker'
        } as msgPayload
    );
}

// 发送消息到 主web 界面
export function sendMainWeb(channel: string, arg: any) {
    ipc.send(ChannelMainWeb, { channel: channel, data: arg, datetime: new Date(), origin: 'worker' });
}


// 发送消息到 Worker（自己发送给自己）
export function sendWorker(channel: string, arg: any) {
    ipc.send(ChannelWorker, { channel: channel, data: arg, datetime: new Date(), origin: 'worker' });
}

// 注册Worker 的监听器, 需要发送至 Worker 的消息需要注册
export function regWorkerListener(callback: (param: msgPayload) => void) {
    ipc.on('worker:message', (event, arg) => {
        callback({
            channel: arg.channel,
            data: arg.data,
            datetime: arg.datetime,
            origin: arg.origin,
        });
    })
}