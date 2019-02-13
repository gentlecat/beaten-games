const path = require('path');

module.exports = {
  entry: ['@babel/polyfill', path.resolve(__dirname, 'src', 'index')],
  output: {
    path: path.resolve(__dirname, 'static', 'js'),
    filename: 'app.bundle.js',
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json'],
    modules: [path.resolve(__dirname, 'src'), 'node_modules'],
  },
  mode: 'development',
  watch: true,
  watchOptions: {
    ignored: ['static', 'node_modules'],
  },

  module: {
    rules: [
      {
        test: /\.(ts|js)x?$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
      },
    ],
  },
};
