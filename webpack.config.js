const path = require("path");

module.exports = {
  entry: {
    home: path.resolve(__dirname, "web/home.ts"),
    game: path.resolve(__dirname, "web/game.ts"),
  },

  output: {
    path: path.resolve(__dirname, "static"),
    filename: "[name].js",
    assetModuleFilename: "assets/[name][ext]",
  },

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
        test: /\.css|\.s(c|a)ss$/,
        use: ["lit-scss-loader", "extract-loader", "css-loader", "sass-loader"],
      },
      {
        test: /\.pdf/,
        type: "asset/resource",
      },
    ],
  },
};
