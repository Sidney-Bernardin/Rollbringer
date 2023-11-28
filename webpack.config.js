const path = require("path");
const sveltePreprocess = require("svelte-preprocess");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = {
  entry: {
    main: path.resolve(__dirname, "web/main.js"),
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
    alias: { svelte: path.resolve("node_modules", "svelte/src/runtime") },
    extensions: [".tsx", ".ts", ".js", ".svelte"],
    mainFields: ["svelte", "browser", "module", "main"],
    conditionNames: ["svelte", "browser", "import"],
  },

  plugins: [
    new MiniCssExtractPlugin({ filename: "[name].css" }),
    new CopyPlugin({ patterns: [{ from: "web/assets", to: "assets" }] }),
  ],

  module: {
    rules: [
      {
        test: /\.(html|svelte)$/,
        use: {
          loader: "svelte-loader",
          options: {
            preprocess: sveltePreprocess(),
            emitCss: true,
          },
        },
      },
      {
        // required to prevent errors from Svelte on Webpack 5+, omit on Webpack 4
        test: /node_modules\/svelte\/.*\.mjs$/,
        resolve: { fullySpecified: false },
      },
      {
        test: /\.(css|scss|sass)$/,
        use: [
          MiniCssExtractPlugin.loader,
          { loader: "css-loader", options: { url: false } }, // https://github.com/peerigon/extract-loader/issues/102#issuecomment-865339845
          "sass-loader",
        ],
      },
    ],
  },
};
