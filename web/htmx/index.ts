window.htmx = require("htmx.org");
import "htmx.org/dist/ext/ws.js";

document.body.addEventListener("htmx:wsConfigSend", (e: CustomEvent) => {
    if (e.detail.parameters["OPERATION"] !== "UPDATE_PDF_FIELD") return;

    for (const k of Object.keys(e.detail.parameters)) {
        if (!k.startsWith("field_value_")) continue;

        e.detail.parameters["field_value"] = e.detail.parameters[k];
        delete e.detail.parameters[k];
        break;
    }
});
