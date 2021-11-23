const MiniCssExtractPlugin = require('mini-css-extract-plugin');

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
        transpileOnly: true,
      },
    },
  },
  // Add support for vue modules
  {
    test: /\.vue$/,
    use: 'vue-loader',
  },
  {
    test: /\.(png|svg|jpg|jpeg|gif)$/i,
    type: 'asset/resource',
  },
  {
    test: /\.(woff|woff2|eot|ttf|otf)$/i,
    type: 'asset/resource',
  },
  // 它会应用到普通的 `.js` 文件
  // 以及 `.vue` 文件中的 `<script>` 块
  // {
  //   test: /\.js$/,
  //   loader: 'babel-loader'
  // },
  // SASS and CSS files from Vue Single File Components:
  {
    test: /\.s[ac]ss$/i,
    use: [
      MiniCssExtractPlugin.loader,
      // { loader: 'vue-style-loader' },
      { loader: 'css-loader' },
      { loader: 'sass-loader' },
    ],
  },
  // normal css files
  {
    test: /\.css$/,
    use: [
      MiniCssExtractPlugin.loader,
      {
        loader: 'css-loader',
        options: {
          importLoaders: 1,
          modules: true,
          sourceMap: true,
        },
      },
    ],
  },
  // markdown
  {
    test: /\.md$/,
    use: [
      {
        loader: 'html-loader',
      },
      {
        loader: 'markdown-loader',
        options: {
          /* your options here */
        },
      },
    ],
  },
  {
    test: /\.worker\.js$/,
    use: [
      {
        loader: 'worker-loader',
        options: { inline: true, fallback: false },
      },
    ],
  },
];
