// import Mdict from 'js-mdict'
// import os from 'os'
// const path = os.path

self.addEventListener('message', function (e) {
  console.log('worker.js')
  // console.log(__static)
  // const medict = new Mdict(path.join(__dirname, '../../static', '/dicts/oale8.mdx'))

  // let def = medict.lookup('hello')
  setTimeout(() => {
    self.postMessage('message from worker: ')
  }, 10000)
}, false)
