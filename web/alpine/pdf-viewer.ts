import PDFContainer from "../utils/pdf-viewer";

window.alpine.data("pdfViewer", (pdfURL: string) => ({
  currentPage: 1,
  pdfViewer: null,

  init() {
    const viewerContainer: HTMLDivElement =
      this.$root.querySelector(".pdf-viewer form");

    this.pdfViewer = new PDFContainer(pdfURL, viewerContainer);

    this.$watch(
      "currentPage",
      async (newVal: string) => await this.pdfViewer.renderPage(newVal),
    );
  },
}));
