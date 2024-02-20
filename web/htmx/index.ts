window.htmx = require("htmx.org");
import "/static/ws.js";

// TODO: Improve error handling.
document.body.addEventListener("htmx:responseError", (e: CustomEvent) => {
    alert(e.detail.xhr.response);
});

document.body.addEventListener("htmx:wsConfigSend", (e: CustomEvent) => {
    const params = e.detail.parameters;

    switch (e.detail.parameters["TYPE"]) {
        case "UPDATE_PDF_PAGE":
            params["pdf_fields"] = {};

            // Moves the message's PDF fields into the nested-object, pdf_fields.
            for (const [k, v] of Object.entries(params)) {
                if (k.startsWith("PDF_")) {
                    params["pdf_fields"][k] = v;
                    delete params[k];
                }
            }

            break;

        case "ROLL":
            if (!(params["die_expressions"] instanceof Array)) {
                params["die_expressions"] = [params["die_expressions"]];
            }

            (params["die_expressions"] as string[]).forEach((dieExpr, idx) => {
                params["die_expressions"][idx] = parseInt(dieExpr);
            });

            break;
    }
});
