const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const { VueLoaderPlugin } = require('vue-loader');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');


module.exports = [
  new ForkTsCheckerWebpackPlugin(),
  new VueLoaderPlugin(),
  new MiniCssExtractPlugin({
      filename:"css~[name]-[chunkhash:12].css"
  }),
  new BundleAnalyzerPlugin({
    analyzerMode: 'static',
    reportFilename: 'electron-webpack-analysis-report.html',
    openAnalyzer: false,
  }),
];
