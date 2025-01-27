const path = require("path");

module.exports = {
    entry: {
        home: path.resolve(__dirname, "web/pages/home.ts"),
    },

    output: {
        path: path.resolve(__dirname, "cmd/static"),
        filename: "[name].js",
    },

    resolve: { extensions: [".ts", ".js"] },

    module: {
        rules: [

            // TypeScript
            {
                test: /\.ts$/,
                use: "ts-loader",
                exclude: /node_modules/,
            },
        ],
    },
};
