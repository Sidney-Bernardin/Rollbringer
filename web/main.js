import "./global.scss";

// Get the page-data from the document.
const script = document.querySelector("[data-page-data]");
const pageData = JSON.parse(script.dataset.pageData);

// Import the page-data's corresponding Svelte page.
const { default: Page } = await import(`./pages/${pageData.name}.svelte`);

export default new Page({
  target: document.body,
  props: pageData,
});
