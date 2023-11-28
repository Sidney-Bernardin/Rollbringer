import "./global.scss";

const tmplData = JSON.parse(document.querySelector("#tmpl-data").textContent);
const { default: Page } = await import(`./pages/${tmplData.name}.svelte`);

export default new Page({ 
  target: document.body,
  props: tmplData,
});
