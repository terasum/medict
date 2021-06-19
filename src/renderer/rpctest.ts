import apis from '../service/service.renderer.register';
for (const fn in apis) {
    if (Object.prototype.hasOwnProperty.call(apis, fn)) {
        console.log(`👨🏻 service renderer process avaliable: ${fn}`);
    }
}
const ret = apis['syncMessage']("myhello");
console.log(`[render-rpc]: syncMessage | ret: ${ret}`);