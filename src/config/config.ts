import { app } from 'electron'
import path from 'path'
import fs from 'fs'

export function getResourceRootPath() {
    const userResourcePath = app.getPath('userData')
    const resourceRootPath = path.resolve(userResourcePath, 'resources', 'cache')
    if (!fs.existsSync(resourceRootPath)) {
        console.log(`create new directory ${resourceRootPath}`)
        fs.mkdirSync(resourceRootPath, {recursive: true})
    }
    return resourceRootPath;
}


