const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const webpack = require('webpack');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');

const FRONTEND_PORT = process.env.FRONTEND_PORT;
const BACKEND_PORT = process.env.BACKEND_PORT;
// const LOCAL_URL = `https://7wt1l5rz-8080.euw.devtunnels.ms/`
const LOCAL_URL = `http://localhost:${BACKEND_PORT}`;

module.exports = (env, argv) => {
  const mode = argv.mode || 'development';
  
  return {
    mode: mode,
    entry: './src/index.tsx',

  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.[contenthash].js',
    clean: true,
    publicPath: '/'
  },
  
  devServer: {
    static: {
      directory: path.join(__dirname, 'dist'),
    },
    port: FRONTEND_PORT,
    hot: true,
    open: true,
    historyApiFallback: true,
    compress: true,
    client: {
      overlay: {
        errors: true,
        warnings: false,
      },
    },
    proxy: {
      '**': {
        target: LOCAL_URL,
        changeOrigin: true,
        secure: false,
        bypass: function(req) {
          if (req.url.startsWith('/static/') || 
              req.url.startsWith('/dist/') ||
              req.url === '/' ||
              req.url.match(/\.(js|css|png|jpg|jpeg|gif|svg|ico|json)$/)) {
            return req.url;
          }
        }
      }
    },
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
    }
  },
  
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              '@babel/preset-env',
              ['@babel/preset-react', { runtime: 'automatic' }],
              '@babel/preset-typescript'
            ]
          }
        }
      },
      {
        test: /\.css$/,
        exclude: /\.module\.css$/,
        use: ['style-loader', 'css-loader']
      },
      {
        test: /\.module\.css$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              modules: {
                localIdentName: '[name]__[local]--[hash:base64:5]'
              }
            }
          }
        ]
      },
      {
        test: /\.(png|svg|jpg|jpeg|gif)$/i,
        type: 'asset/resource'
      }
    ]
  },
  
  plugins: [
    new HtmlWebpackPlugin({
      template: './public/index.html',
      filename: 'index.html',
      inject: 'body'
    }),
    new Dotenv({
      path: './.env'
    }),
    new ForkTsCheckerWebpackPlugin({
      typescript: {
        diagnosticOptions: {
          semantic: true,
          syntactic: true,
        },
      },
    }),
    new webpack.DefinePlugin({
      'process.env.ENV': JSON.stringify(process.env.ENV || 'development'),
      'process.env.IS_PROD': JSON.stringify(process.env.ENV === 'production'),
      'process.env.NODE_ENV': JSON.stringify(mode),
      'IS_LOCAL_HOST': JSON.stringify(mode === 'development'),
      'IS_DEV': JSON.stringify(mode === 'development')
    })
  ],
  
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx'],
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  
  devtool: mode === 'development' ? 'eval-source-map' : 'source-map',
  
  stats: {
    colors: true,
    modules: false,
    reasons: true,
    errorDetails: true,
    performance: true,
    warnings: true,
    errors: true,
    errorStack: true,
    moduleTrace: true
  }
  };
};