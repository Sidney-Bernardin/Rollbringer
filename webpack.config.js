const path = require("path");

module.exports = {
  entry: {
    play: path.resolve(__dirname, "web/play.ts"),

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
        use: [
          "style-loader",
          {
            loader: "css-loader",
            options: { url: false },
          },
          "sass-loader",
        ],
      },
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
