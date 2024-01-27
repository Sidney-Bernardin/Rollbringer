import PDFContainer from "../utils/pdf-viewer";

window.alpine.data("pdfViewer", (pdfURL: string) => ({
  pdfViewer: null,

  init() {
    const initPageForm: HTMLFormElement = this.$root.querySelector(
      ".pdf-viewer__init-page-form",
    );

    const viewerContainer: HTMLDivElement = this.$root.querySelector(
      ".pdf-viewer__container",
    );

    this.pdfViewer = new PDFContainer(pdfURL, initPageForm, viewerContainer);

    this.$watch(
      "currentPage",
      async (newVal: string) => await this.pdfViewer.renderPage(newVal),
    );
  },
}));
