<script lang="ts">
  import { onMount } from "svelte";
  import { PDFDocumentProxy, getDocument } from "pdfjs-dist";
  import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";

  export let pdfURL: string;

  let viewerContainer: HTMLDivElement;

  let pdfDoc: PDFDocumentProxy;
  let pageView: PDFPageView;
  let currentPage: number;

  onMount(async () => {
    pdfDoc = await getDocument(pdfURL).promise;
    const firstPage = await pdfDoc.getPage(1);

    pageView = new PDFPageView({
      id: 1,
      container: viewerContainer,
      eventBus: new EventBus(),
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });
  });
</script>

<!-- markup (zero or more items) goes here -->

<style>
  /* your styles go here */
</style>
