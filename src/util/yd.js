import axios from 'axios'

const API = 'YouDaoCV'
const API_KEY = '659600698'

const englishReg = /^[a-z|A-Z]+$/
const chineseReg = /[\u4e00-\u9fff]+$/

const onlineRefs = [
  {
    lang: englishReg,
    url: 'http://www.ldoceonline.com/search/?q={0}',
    name: '朗文在线'
  },
  {
    lang: englishReg,
    url: 'http://dictionary.reference.com/browse/{0}',
    name: 'reference.com'
  },
  {
    lang: englishReg,
    url: 'http://www.urbandictionary.com/define.php?term={0}',
    name: 'urbandictionary.com'
  },
  {
    lang: chineseReg,
    url: 'http://www.zdic.net/sousuo/?q={0}',
    name: 'zdic在线辞海'
  }
]

const youdaoUrl = 'http://fanyi.youdao.com/openapi.do?keyfrom={0}&key={1}&type=data&doctype=json&version=1.2&q={2}'

/**
 * YouDao dictionary wrapper class
 * @param {object} options : the youdao dict api options, usually be {API:'', API_KEY: ''}
 */
function YouDaoDict (options) {
  let config = {}
  if (!options) {
    Object.assign(config, options)
  }
  this.API = config.API || API
  this.API_KEY = config.API || API_KEY
}

YouDaoDict.prototype.lookup = async function (word) {
  let queryUrl = youdaoUrl
    .replace('{0}', this.API)
    .replace('{1}', this.API_KEY)
    .replace('{2}', word)
  try {
    let response = await axios.get(queryUrl, {
      headers: {
        'Access-Control-Allow-Origin': '*'
      }
    })
    // console.log(response.data)
    return response.data
  } catch (e) {
    console.log(e)
    return '* NOT FOUND *'
  }
}

YouDaoDict.prototype.formatHTML = function (def) {
  let querySect = `<section class="yd-query">
    <span class="yd-query">${def.query}</span>
    <span class="yd-phonetic">${def.basic['phonetic']}</span>
    <span class="yd-speech">
      <audio id="phonetic" class="yd-audio">
      <source src="${def.basic['speech']}" type="audio/mpeg">
      </audio>
    </span>`

  if (def.basic['uk-phonetic']) {
    querySect += `
    <span class="yd-uk-phonetic-label">UK</span>
    <span class="yd-uk-phonetic">${def.basic['uk-phonetic']}</span>
    <span class="yd-uk-speech">
      <audio id="uk-phonetic" class="yd-audio">
        <source src="${def.basic['uk-speech']}" type="audio/mpeg">
      </audio>
      </span>`
  }

  if (def.basic['us-phonetic']) {
    querySect += `
    <span class="yd-us-phonetic-label">US</span>
    <span class="yd-us-phonetic">${def.basic['us-phonetic']}</span>
    <span class="yd-us-speech">
      <audio id="us-phonetic" class="yd-audio">
      <source src="${def.basic['us-speech']}" type="audio/mpeg">
      </audio>
    </span>
    `
  }
  querySect += `</section>`

  let explainSect = `<section class="yd-explain">
      <ul class="yd-wordexps">`
  for (let i = 0; i < def.basic.explains.length; i++) {
    explainSect += `<li class="yd-expitem">${def.basic.explains[i]}</li>`
  }
  explainSect += `</ul>
    </section>`

  let webrefSect = `  <section class="yd-webref">
    <ul class="yd-webrefs">`
  for (let i = 0; i < def.web.length; i++) {
    webrefSect += `<li class="yd-webrefitem">
      <span class="yd-webrefitem-key">${def.web[i].key}</span>
      <span class="yd-webrefitem-value">${def.web[i].value.join('; ')}</span>
      </li>`
  }
  webrefSect += `</ul>
  </section>`

  let onlineResSect = ` <section class="yd-onlineres">
  <ul class="yd-online-ress">`

  onlineRefs.forEach(element => {
    if (element.lang.test(def.query)) {
      onlineResSect += `<li class="yd-onlineresitem">
      <a class="yd-onlinres-item-link" href="${element.url.replace('{0}', def.query)}" target="_blank"> ${element.name}</a>
      </li>`
    }
  })

  onlineResSect += `</ul>
</section>`

  return querySect +
    explainSect +
    webrefSect +
    onlineResSect
}

// const youdao = new YouDaoDict()
// youdao.lookup('section').then((def) => {
//   console.log(youdao.formatHTML(def))
// })

// console.log(chineseReg.test('中文'))
// console.log(chineseReg.test('english'))
// console.log(englishReg.test('english'))
// console.log(englishReg.test('English'))
// console.log(englishReg.test('中文'))

export default YouDaoDict
