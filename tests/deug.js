var str = `<link rel="stylesheet" type="text/css" href="oalecd8e.css"><script src="jquery.js" charset="utf-8" type="text/javascript" \ language="javascript"></script><script src="oalecd8e.js" charset="utf-8" type="text/javascript" language="javascript"></script><span 
id="weakness_e" name="weakness" idm_id="000040789" class="entry"><span class="h-g"><span class="top-g"><span class="h">weak·ness</span><span class="oalecd8e_show_all"><em></em></span> <span class="z"> <span class="symbols-coresym">★</span> </span><span class="ei-g"><span class="z_ei-g">/</span><a class="fayin" href="sound://uk/weakness__gb_1.mp3"><span class="phon-gb">ˈwiːknəs</span><img src="uk_pron.png" class="fayin"/></a><span class="z">; <span class="z_phon-us">NAmE</span></span><a class="fayin" href="sound://us/weakness__us_1.mp3"><span class="phon-usgb">ˈwiːknəs</span><img src="us_pron.png" class="fayin"/></a><span class="z_ei-g">/</span></span><span class="block-g"><span class="pos-g"> <span class="pos">noun</span> </span></span></span><span id="weakness_ng_1" class="n-g"><span class="z_n">1. </span><span class="symbols-small_coresym">★</span> <span class="gr"><span class="z_gr_br"> [</span>U<span class="z_gr_br">] </span></span> <span class="def-g"><span class="d oalecd8e_switch_lang switch_children">lack of strength, power or determination <span class="oalecd8e_chn">软弱；虚弱；疲软；衰弱；懦弱</span></span></span><span id="weakness_xg_1" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">The sudden weakness in her legs made her stumble. </span><span class="oalecd8e_chn">她突然两腿发软踉跄了一下。</span></span><span id="weakness_xg_2" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">the weakness of the dollar against the pound </span><span class="oalecd8e_chn">美元对英镑的疲软</span></span><span id="weakness_xg_3" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">He thought that crying was a <span class="cl">sign of weakness.</span> </span><span class="oalecd8e_chn">他认为哭是懦弱的表现。</span></span><span class="xr-g"> <span class="symbols-oppsym"><a href="#O8T">OPP</a></span> <span id="weakness_xr_1" href="strength_e" class="xr"><span class="Ref"> <span class="xh"> <a href="entry://strength">strength</a></span> </span></span></span></span><span id="weakness_ng_2" class="n-g"><span class="z_n">2. </span><span class="symbols-small_coresym">★</span> <span class="gr"><span class="z_gr_br"> [</span>C<span class="z_gr_br">] </span></span> <span class="def-g"><span class="d oalecd8e_switch_lang switch_children">a weak point in a system, sbs character, etc. <span class="oalecd8e_chn">（系统、性格等的）弱点，缺点，不足</span></span></span><span id="weakness_xg_4" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">Its important to know your own <span class="cl">strengths and weaknesses.</span> </span><span class="oalecd8e_chn">了解自己的优缺点很重要。</span></span><span id="weakness_xg_5" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">Can you spot the weakness in her argument? </span><span class="oalecd8e_chn">你能指出她论点中的不足之处吗？</span></span><span class="xr-g"> <span class="symbols-oppsym"><a href="#O8T">OPP</a></span> <span id="weakness_xr_2" href="strength_e" class="xr"><span class="Ref"> <span class="xh"> <a href="entry://strength">strength</a></span> </span></span></span></span><span id="weakness_ng_3" class="n-g"><span class="z_n">3. </span><span class="gr"><span class="z_gr_br"> [</span>C</span><span class="z">, </span><span class="gr">usually sing.<span class="z_gr_br">] </span></span><span id="weakness_cf_1" class="cf"> <span class="swung-dash">weakness</span>~ (for sth/sb) </span><span class="def-g"><span class="d oalecd8e_switch_lang switch_children">difficulty in resisting sth/sb that you like very much <span class="oalecd8e_chn">（对人或事物的）迷恋，无法抗拒</span></span></span><span id="weakness_xg_6" class="x-g"><span class="symbols-xsym"></span><span class="x oalecd8e_switch_lang switch_siblings">He has a weakness for chocolate. </span><span class="oalecd8e_chn">他爱吃巧克力。</span></span></span><span class="infl"><span class="inflection">weakness</span> <span class="inflection">weaknesses</span> </span></span><span class="pracpron"><span class="pron-g"><span class="wd">weak·ness</span> <span class="ei-g"><span class="z_ei-g">/</span><a class="fayin" href="sound://uk/weakness__gb_1.mp3"><span class="phon-gb">ˈwiːknəs</span><img src="uk_pron.png" class="fayin"/></a><span class="z">; <span class="z_phon-us">NAmE</span></span><a class="fayin" href="sound://us/weakness__us_1.mp3"><span class="phon-usgb">ˈwiːknəs</span><img src="us_pron.png" class="fayin"/></a>
<span class="z_ei-g">/</span></span></span></span></span>`

const IMAGE_REG = /src=\"((\S+)\.(png|jpg|gif|jpeg|svg))\"/gi
const IMAGE_REG_IDX = 1

const SOUND_REG = /sound\:\/\/((\S+)\.(mp3))/gi
const SOUND_REG_IDX = 0

const CSS_REG = /href=\"((\S+)\.css)\"/
const CSS_REG_IDX = 1

const JS_REG = /src=\"((\S+)\.js)\"/
const JS_REG_IDX = 1

let matches = str.matchAll(IMAGE_REG)
for (const match of matches) {
  console.log(
    `Found ${match[IMAGE_REG_IDX]} start=${match.index} end=${
      match.index + match[IMAGE_REG_IDX].length
    }.`,
  )
}

console.log('--------------------')
matches = str.matchAll(SOUND_REG)
for (const match of matches) {
  console.log(
    `Found ${match[SOUND_REG_IDX]} start=${match.index} end=${
      match.index + match[SOUND_REG_IDX].length
    }.`,
  )
}

console.log('--------------------')
matches = str.matchAll(CSS_REG)
for (const match of matches) {
  console.log(
    `Found ${match[CSS_REG_IDX]} start=${match.index} end=${
      match.index + match[CSS_REG_IDX].length
    }.`,
  )
}

console.log('--------------------')

matches = str.matchAll(JS_REG)
for (const match of matches) {
  console.log(
    `Found ${match[JS_REG_IDX]} start=${match.index} end=${
      match.index + match[JS_REG_IDX].length
    }.`,
  )
}
