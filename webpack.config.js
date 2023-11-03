const path = require("path");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: {
    home: path.resolve(__dirname, "web/pages/home.ts"),
    game: path.resolve(__dirname, "web/pages/game.ts"),
    "pdf.worker": path.resolve(
      __dirname,
      "node_modules/pdfjs-dist/build/pdf.worker.js",
    ),
  },

  output: {
    path: path.resolve(__dirname, "static"),
    filename: "[name].js",
    assetModuleFilename: "assets/[name][ext]",
  },

  plugins: [
    new CopyPlugin({
      patterns: [{ from: "web/assets", to: "assets" }],
    }),
  ],

  resolve: {
    extensions: [".tsx", ".ts", ".js"],
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
      {
        test: /\.pdf/,
        type: "asset/resource",
      },
      {
        test: /\.css|\.s(c|a)ss$/,
        use: [
          "lit-scss-loader",
          "extract-loader",
          {
            loader: "css-loader",
            options: { url: false }, // https://github.com/peerigon/extract-loader/issues/102#issuecomment-865339845
          },
          "sass-loader",
        ],
      },
    ],
  },
};
