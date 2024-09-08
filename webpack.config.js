const path = require('path');

module.exports = {
    entry: '/static/js/dashboard/main.js',
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'static/js/dashboard'),
    },
    mode: 'development',
};
