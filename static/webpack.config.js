var webpack = require('webpack');
var path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');


var BUILD_DIR = path.resolve(__dirname, '');
var APP_DIR = path.resolve(__dirname, '');

const HtmlWebpackPluginConfig = new HtmlWebpackPlugin({
	template: 'index.template.ejs',
	filename: 'index.html',
	inject: 'body'
})



var config = {
	entry: APP_DIR + '/index.jsx',
	output: {
		path: BUILD_DIR,
		filename: 'bundle.js'
	}
	,
	module: {
		rules: [{
			test: /\.jsx?/,
			include: APP_DIR,
			exclude: /node_modules/,
			use: [
				"babel"
			]
		}]
	}
	, resolveLoader: {
		moduleExtensions: ["-loader"]
	}
	,
	plugins: [HtmlWebpackPluginConfig]
};


module.exports = config;
