// https://www.postcss.parts/
const autoprefixer = require('autoprefixer')
const postcss_pxtorem = require('postcss-pxtorem')

module.exports = {
  plugins: [autoprefixer, postcss_pxtorem],
}
