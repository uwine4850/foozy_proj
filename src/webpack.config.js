const miniCss = require('mini-css-extract-plugin');
module.exports = {
    // mode: "development",
    entry: './static/js/index.js',
    output: {
        filename: 'bundle.js',
        path: '/usr/app'
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [
                    miniCss.loader,
                    'css-loader',
                ],
            },
            {
            test:/\.(s*)css$/,
            use: [
                miniCss.loader,
                'css-loader',
                'sass-loader',
            ]
        }]
    },
    plugins: [
        new miniCss({
            filename: './static/css/style.css',
        }),
    ],
    resolve: {
        extensions: ['', '.js', '.es6', '.jsx']
    }
};