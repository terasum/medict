const rules = require('./webpack.rules');
const plugins = require('./webpack.plugins');
const path = require('path');

rules.push({
  test: /\.scss$/,
  use: [
    { loader: 'vue-style-loader' },
    { loader: 'style-loader' },
    { loader: 'css-loader' },
    {
      // Run postcss actions
      loader: 'postcss-loader',
      options: {
        // `postcssOptions` is needed for postcss 8.x;
        // if you use postcss 7.x skip the key
        postcssOptions: {
          // postcss plugins, can be exported to postcss.config.js
          plugins: function() {
            return [require('autoprefixer')];
          },
        },
      },
    },
    {
      // compiles Sass to CSS
      loader: 'sass-loader',
    },
  ],
});

// normal css files
rules.push({
  test: /\.css$/,
  use: ['vue-style-loader', 'style-loader', 'css-loader'],
});

module.exports = {
  module: {
    rules,
  },
  plugins: plugins,
  resolve: {
    extensions: ['.js', '.ts', '.jsx', '.tsx', '.css', '.scss', '.vue'],
    alias: {
      vue$: 'vue/dist/vue.esm.js', // 用 webpack 1 时需用 'vue/dist/vue.common.js'
      // '@': path.resolve(__dirname, 'node_modules/'),
      // '~': path.resolve(__dirname, 'src/renderer/'),
    },
  },
};
