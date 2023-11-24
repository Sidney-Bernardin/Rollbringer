<script context="module">
  GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";
</script>

<script lang="ts">
  import { onMount } from "svelte";
  import Panzoom from "panzoom";
  import {
    GlobalWorkerOptions,
    PDFDocumentProxy,
    getDocument,
  } from "pdfjs-dist";
  import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";

  let viewerContainer: HTMLDivElement;

  export let pdfURL: string;
  let pdfDoc: PDFDocumentProxy;
  let pageView: PDFPageView;
  export let currentPage: number;

  $: currentPage,
    (async function () {
      if (!pageView || !pdfDoc) return;
      pageView.setPdfPage(await pdfDoc.getPage(currentPage));
      pageView.draw();
    })();

  onMount(async () => {
    pdfDoc = await getDocument(pdfURL).promise;
    const firstPage = await pdfDoc.getPage(currentPage);

    pageView = new PDFPageView({
      id: 1,
      container: viewerContainer,
      eventBus: new EventBus(),
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });

    Panzoom(viewerContainer, {
      bounds: true,
      minZoom: 0.25,
      maxZoom: 2,
      smoothScroll: false,
      filterKey: () => true,
    });
  });
</script>

<div class="pdf-viewer">
  <div class="viewer-container" bind:this={viewerContainer}>
    <div class="viewer" />
  </div>
</div>

<style lang="scss">
  @import "~pdfjs-dist/web/pdf_viewer.css";

  .pdf-viewer {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;

    .viewer-container {
      position: absolute;

      .textLayer {
        user-select: none;
      }
    }
  }
</style>
