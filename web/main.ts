require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);
import "./styles/pages/play/index.scss";

import { Alpine } from "alpinejs";
import { HtmxExtension } from "htmx.org";

declare global {
    interface Window {
        alpine: Alpine;
        htmx: HtmxExtension;
    }
}

import "./alpine";
window.htmx = require('htmx.org');

window.addEventListener("htmx:wsConfigSend", (e: CustomEvent) => {
    switch (e.detail.parameters["EVENT"]) {
        case "UPDATE_PDF_PAGE_REQUEST":
            e.detail.parameters["field_name"] = e.detail.headers["HX-Trigger"];
            for (const k of Object.keys(e.detail.parameters)) {
                if (!k.startsWith("field_value_")) continue;

                e.detail.parameters["field_value"] = e.detail.parameters[k]
                delete e.detail.parameters[k];
                break;
            }
            break;
    }
});