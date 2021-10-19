const path = require('path');

const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const { VueLoaderPlugin } = require('vue-loader');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

const assets = ['docs', 'css']; // asset directories

const plugins = [
  new ForkTsCheckerWebpackPlugin(),
  new VueLoaderPlugin(),
  new MiniCssExtractPlugin({
    filename: 'css~[name]-[chunkhash:12].css',
  }),
  new CopyWebpackPlugin({
    patterns: assets.map((asset) => {
      return {
        from: path.resolve(__dirname, 'src/renderer/assets', asset),
        to: path.resolve(__dirname, '.webpack/renderer/assets', asset),
      };
    }),
  }),
  new BundleAnalyzerPlugin({
    analyzerMode: 'static',
    reportFilename: 'electron-webpack-analysis-report.html',
    openAnalyzer: false,
  }),
];

module.exports = plugins;
