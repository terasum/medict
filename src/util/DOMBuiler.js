import { JSDOM } from 'jsdom'

/**
 * wrapper content as a html
 * @param {string} content wrappered content
 */
function wrapper (content) {
  return `<!DOCTYPE html>
  <html lang="zh-cn">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>medict-content</title>
  </head> 
  <body>
  <div class="medict-wrapper">
  ` +
  content +
  `</div>
  </body>
  </html>`
}

function DOMBuilder (def) {
  this.dom = new JSDOM(wrapper(def))
  /**
 * add script connect
 * @param {*} scriptContent JavaScript raw content
 */
  this.addScriptContent = (scriptContent) => {
    let bodytag = this.dom.window.document.body
    let scriptNode = this.dom.window.document.createElement('SCRIPT')
    scriptNode.textContent = scriptContent
    bodytag.appendChild(scriptNode)
  }

  /**
 * add script path
 * @param {string} scriptPath script load url
 */
  this.addScript = (scriptPath) => {
    let bodytag = this.dom.window.document.body
    let scriptNode = this.dom.window.document.createElement('SCRIPT')
    scriptNode.setAttribute('src', scriptPath)
    scriptNode.setAttribute('data-medict', 'true')
    bodytag.appendChild(scriptNode)
  }

  /**
 * add style content into dom
 * @param {string} styleContent style raw content
 */
  this.addStyleContent = (styleContent) => {
    let headertag = this.dom.window.document.head
    let styleNode = this.dom.window.document.createElement('STYLE')
    headertag.appendChild(styleNode)
    styleNode.textContent = styleContent
  }

  /**
 * add style script by url
 * @param {string} stylePath style path
 */
  this.addStyle = (stylePath) => {
    let bodytag = this.dom.window.document.body
    let linkNode = this.dom.window.document.createElement('LINK')
    linkNode.setAttribute('src', stylePath)
    linkNode.setAttribute('data-medict', 'true')
    bodytag.appendChild(linkNode)
  }

  /**
 * _normal append script
 */
  this.normal = () => {
    const textContent = ``
    this.addScriptContent(textContent)
  }

  /**
 * serialize
 */
  this.serialize = () => {
    console.log(this)
    this.normal()
    return this.dom.serialize()
  }
}

export default DOMBuilder
