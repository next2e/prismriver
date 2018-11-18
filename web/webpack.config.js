const CopyWebpackPlugin = require('copy-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const { resolve } = require('path')

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
        test: /\.ts?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      }
    ]
  },
  resolve: {
    extensions: ['.css', '.ts']
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
    })
  ]
}