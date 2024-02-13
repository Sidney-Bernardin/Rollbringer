window.htmx = require("htmx.org");
import "/static/ws.js";

// TODO: Improve error handling.
document.body.addEventListener("htmx:responseError", (e: CustomEvent) => {
    alert(e.detail.xhr.response);
});

document.body.addEventListener("htmx:wsConfigSend", (e: CustomEvent) => {
    if (e.detail.parameters["TYPE"] === "UPDATE_PDF_PAGE") {
        e.detail.parameters["pdf_fields"] = {};

        // Moves the message's PDF fields into the nested-object, pdf_fields.
        for (const [k, v] of Object.entries(e.detail.parameters)) {
            if (k.startsWith("PDF_")) {
                e.detail.parameters["pdf_fields"][k] = v;
                delete e.detail.parameters[k];
            }
        }
    }
});
