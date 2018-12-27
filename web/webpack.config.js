const CopyWebpackPlugin = require('copy-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const { resolve } = require('path')
const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
  entry: {
    appJs: './src/js/app.ts'
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader
          },
          'css-loader'
        ]
      },
      {
        test: /\.(woff(2)?|ttf|eot|svg)(\?v=\d+\.\d+\.\d+)?$/,
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[ext]',
            outputPath: './fonts',
            publicPath: '../fonts'
          }
        }]
      },
      {
        test: /\.html$/,
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[ext]',
            outputPath: './'
          }
        }]
      },
      {
        loader: 'ts-loader',
        options: {
          appendTsSuffixTo: [/\.vue$/],
        },
        test: /\.ts?$/,
        exclude: /node_modules/
      },
      {
        loader: 'vue-loader',
        test: /\.vue$/
      }
    ]
  },
  resolve: {
    extensions: ['.css', '.js', '.ts']
  },
  output: {
    filename: 'js/app.js',
    path: resolve(__dirname, 'dist')
  },
  plugins: [
    new CopyWebpackPlugin([
      { from: './src/img', to: 'img' },
      { from: './src/index.html', to: 'index.html' }
    ]),
    new MiniCssExtractPlugin({
      filename: 'css/app.css'
    }),
    new VueLoaderPlugin()
  ]
}