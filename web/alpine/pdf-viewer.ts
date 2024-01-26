import PDFContainer from "../utils/pdf-viewer";

window.alpine.data("pdfViewer", (pdfURL: string) => ({
  pdfViewer: null,

  init() {
    const container: HTMLDivElement = this.$root.querySelector(
      ".pdf-viewer__container",
    );

    this.pdfViewer = new PDFContainer(pdfURL, container);

    this.$watch("currentPage", async (newVal: string) => await this.pdfViewer.renderPage(newVal));
  },
}));
