import DOMBuilder from '../util/dombuiler'

class Query {
  constructor (name, def, config) {
    this.name = name
    this.config = config || {}
    this.content = new DOMBuilder(def)
    this.decorate = false
  }

  wrapper () {
    if (!this.decorate) {
      throw new Error('should invoke decorate first')
    }
  }

  decorate () {
    this.decorated = true
    this.content.addStyleContent(`
      html {
        height: 100%;
        width: 557px !important;
      }
      body {
      height:100%;
      width: 527px !important;
      cursor: default;
      margin: 0;
      padding: 5px;
      padding-bottom: 40px;
      overflow-y: scroll;
      }
      .medict-wrapper{
        padding: 2px;
        width: 100%;
        height:100%;
        padding-bottom:50px;
      }
    `)
    // this.content.addScript('https://code.jquery.com/jquery-2.2.4.js')
    this.content.addScript('file://' + __static + '/scripts/jquery-2.2.4.js')
    this.content.addScriptContent(`
     window.jQuery = window.$ = module.exports;
     $('a').on('click', function(event) {
       console.log(event.target);
       const url = event.target.getAttribute('href')
       console.log(url);
       if(url && url.startsWith('http')){
        console.log('open in browser');
       }else{
        event.preventDefault();
       }
     })
    `)
  }

  serialize () {
    return `data:text/html, ` + this.content.serialize()
  }
}
export default Query
