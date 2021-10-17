const rules = require('./webpack.rules');
const plugins = require('./webpack.plugins');

const TerserJSPlugin = require('terser-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');



module.exports = {
  module: {
    rules,
  },
  plugins: plugins,
  optimization: {
    minimize: true,
    minimizer: [
      new TerserJSPlugin({
        parallel: true,
      }),
      new OptimizeCssAssetsPlugin({
        cssProcessorOptions: {
          map: {
            inline: false,
            annotation: true,
          },
        },
      }),
    ],
    splitChunks: {
      chunks: 'all', //表示显示块的范围，有三个可选值：initial(初始块)、async(按需加载块)、all(全部块)(default=all);
      minChunks: 1, //在分割之前模块的被引用次数(default=1)
      minSize: 30000, //代码块的最小尺寸(default=30000)
      maxSize: 50000,
      maxAsyncRequests: 10, //按需加载最大并行请求数量(default=5)
      maxInitialRequests: 5, //一个入口的最大并行请求数量(default=3)
      cacheGroups: {
        //可以继承/覆盖上面 splitChunks 中所有的参数值
        vue: {
          test: /[\\/]node_modules[\\/](vue)[\\/]/,
          name: 'vue',
          priority: 2,
        },
        router: {
          test: /[\\/]node_modules[\\/](vue-router)[\\/]/,
          name: 'vue-router',
          priority: 2,
        },
        router: {
          test: /[\\/]node_modules[\\/](buefy)[\\/]/,
          name: 'buefy',
          priority: 3,
        },
        vendors: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          priority: 1,
          reuseExistingChunk: true,
        },
      },
    },
    runtimeChunk: {
      name: 'manifest',
    },
  },
  resolve: {
    extensions: ['.js', '.ts', '.jsx', '.tsx', '.css', '.scss', '.vue'],
    alias: {
      vue$: 'vue/dist/vue.esm.js', // 用 webpack 1 时需用 'vue/dist/vue.common.js'
      // '@': path.resolve(__dirname, 'node_modules/'),
      // '~': path.resolve(__dirname, 'src/renderer/'),
    },
  },
  output: {
    filename: 'js~[name].[chunkhash:12].js',
  },
};
