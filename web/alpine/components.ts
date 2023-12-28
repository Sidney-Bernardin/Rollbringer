import Alpine from "alpinejs";
import PDFContainer from "../utils/pdf-viewer";

Alpine.data("pdfViewer", (pdfURL: string) => ({
  pdfViewer: null,

  init() {
    const container = this.$root.querySelector(".pdf-viewer__container")

    this.pdfViewer = new PDFContainer(pdfURL, container as HTMLDivElement)

    this.$watch("currentTab", async (newVal: string) => await this.pdfViewer.renderPage(newVal));
  },
}));
