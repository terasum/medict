const electron = require('electron')
const path = require('path')
const fs = require('fs')

const userDataPath = (electron.app || electron.remote.app).getPath('userData')

// reference from: https://medium.com/cameron-nokes/how-to-store-user-data-in-electron-3ba6bf66bc1e
class FileStore {
  constructor (relativePath) {
    // Renderer process has to get `app` module via `remote`, whereas the main process can get it directly
    // app.getPath('userData') will return a string of the user's app data directory path.
    // We'll use the `configName` property to set the file name and path.join to bring it all together as a string
    this.path = path.join(userDataPath, relativePath)
  }

  writeRawData (rawData) {
    fs.writeFileSync(this.path, rawData)
  }

  readRawData () {
    return fs.readFileSync(this.path)
  }
}

// expose the class
export default {FileStore, userDataPath}
