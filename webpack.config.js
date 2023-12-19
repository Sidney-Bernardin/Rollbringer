const path = require("path");

module.exports = {
  entry: {
    home: path.resolve(__dirname, "web/home.ts"),
    game: path.resolve(__dirname, "web/game.ts"),

    "pdf.worker": path.resolve(
      __dirname,
      "node_modules/pdfjs-dist/build/pdf.worker.js",
    ),
  },

  output: {
    path: path.resolve(__dirname, "static"),
    filename: "[name].js",
  },

  resolve: {
    extensions: [".ts", ".js"],
  },

  module: {
    rules: [
      {
        test: /\.ts$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
      {
        test: /\.s[ac]ss$/i,
        use: ["style-loader", "css-loader", "sass-loader"],
      },
    ],
  },
};
