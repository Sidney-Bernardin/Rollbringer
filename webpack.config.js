const path = require("path");

module.exports = {
    entry: {
        play: path.resolve(__dirname, "web/play.ts"),

        ws: path.resolve(__dirname, "node_modules/htmx.org/dist/ext/ws.js"),
        "pdf.worker": path.resolve(__dirname, "node_modules/pdfjs-dist/build/pdf.worker.js"),
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

            // SASS
            {
                test: /\.s[ac]ss$/i,
                use: [
                    "style-loader",
                    {
                        loader: "css-loader",
                        options: { url: false },
                    },
                    "sass-loader",
                ],
            },

            // Assets
            {
                test: /\.(png|jpg|gif|pdf)$/,
                use: [
                    {
                        loader: "file-loader",
                        options: {
                            name: "[name].[ext]",
                        },
                    },
                ],
            },
        ],
    },
};
