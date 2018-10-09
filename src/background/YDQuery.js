import Query from './Query'
class YDQuery extends Query {
  constructor (def, config) {
    super('youdao', def, config)
  }

  wrapper () {
    super.decorate()
    // style
    this.content.addStyleContent(`
      body{ background: #fcfcfc;
      font-family: Verdana,Helvetica,Arial,sans-serif;
      }
      
      ul{
        -webkit-padding-start: 20px;
        font-size: 14px;
      }
      li{
        list-style: none;
      }

      /* query part */
      section.yd-query{
      display: block;
      font-family: Verdana,Helvetica,Arial,sans-serif;
      margin-top: 10px;
      }

      span.yd-query{
      display:block;
      width:100%;
      font-size:26px;
      margin-left: 15px;
      margin-bottom: 10px;
      border-bottom: 1px solid #227ee8;
      padding-bottom: 5px;
      }

      span.yd-phonetic {
      background: #eee;
      margin-left:15px;
      font-size: 16px;
      cursor: pointer;
      }
      span.yd-phonetic::before {
        content: "/";
      }
      span.yd-phonetic::after {
        content: "/";
      }

      span.yd-us-phonetic-label{
        font-size:14px;
        color: #888;
        margi-left:7px;
      }
      span.yd-us-phonetic{
      background: #eee;
      color: #999;
      font-size: 14px;
      cursor: pointer;
      }
      span.yd-us-phonetic::before {
        content: "/";
      }
      span.yd-us-phonetic::after {
        content: "/";
      }

      span.yd-uk-phonetic-label{
        font-size:14px;
        color: #888;
        margi-left:7px;
      }
      span.yd-uk-phonetic{
      background: #eee;
      color: #999;
      font-size: 14px;
      cursor: pointer;
      }
      span.yd-uk-phonetic::before {
        content: "/";
      }
      span.yd-uk-phonetic::after {
        content: "/";
      }
      /* web references */
      section.yd-webrefs{
        background: aliceblue;
        padding: 15px 0 15px 20px;
      }
      li.yd-webrefitem{
        list-style: decimal;
        color:#444;
        font-size: 14px;
        margin-left: 20px;
      }

      /* online references part */
      section.yd-webref{
        padding: 10px 0px;
        background-color: aliceblue;
      }
      li.yd-onlineresitem {
        margin-left: 20px;
        list-style: circle;
        color: #666;
      }

      a.yd-onlinres-item-link{
        color: rgb(53, 161, 212);
        font-size: 13px;
        font-family: Verdana,Helvetica,Arial,sans-serif;
      }
    `)
    this.content.addScriptContent(`
      console.log("loading....")
      console.log(jQuery)
      $(document).ready(function(){
        var phonetic = $('.yd-phonetic')
        var phonetic_sp = $('#phonetic')[0]
        if (phonetic && phonetic_sp ){
            phonetic.on('click', function(){
              console.log('clicked phonetic')
              phonetic_sp.play()
            })
        }

        var ukphonetic = $('.yd-uk-phonetic')
        console.log(ukphonetic)
        var ukphonetic_sp = $('#uk-phonetic')[0]
        if (ukphonetic &&  ukphonetic_sp){
            ukphonetic.on('click', function(){
              console.log('clicked uk phonetic')
              ukphonetic_sp.play()
            })
        }
        var usphonetic = $('.yd-us-phonetic')
        console.log(usphonetic)
        var usphonetic_sp = $('#us-phonetic')[0]
        if (usphonetic && usphonetic_sp){
            usphonetic.on('click', function(){
              console.log('clicked us phonetic')
              usphonetic_sp.play()
            })
        }
      })
    `)
  }
}

export default YDQuery
