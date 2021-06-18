module.exports = [
  // Add support for native node modules
  {
    test: /\.node$/,
    use: 'node-loader',
  },
  {
    test: /\.(m?js|node)$/,
    parser: { amd: false },
    use: {
      loader: '@marshallofsound/webpack-asset-relocator-loader',
      options: {
        outputAssetBase: 'native_modules',
      },
    },
  },
  {
    test: /\.tsx?$/,
    exclude: /(node_modules|\.webpack)/,
    use: {
      loader: 'ts-loader',
      options: {
        transpileOnly: true
      }
    }
  },
  // Add support for vue modules
  {
    test: /\.vue$/,
    use: 'vue-loader',
  },
   // 它会应用到普通的 `.js` 文件
   // 以及 `.vue` 文件中的 `<script>` 块
    // {
    //   test: /\.js$/,
    //   loader: 'babel-loader'
    // },
];
