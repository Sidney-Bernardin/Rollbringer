window.htmx = require('htmx.org');
import "htmx.org/dist/ext/ws.js";


window.addEventListener("htmx:responseError", (e: CustomEvent) => {
    let res = JSON.parse(e.detail.xhr.response)
    alert(`${res.type}: ${res.detail}`)
});

window.addEventListener("htmx:configRequest", (e: CustomEvent) => {
    const currentGameID = document.body.dataset.gameId
    const minimizedResponseGameID = e.detail.elt.dataset.minimizedResponseGameId

    if (minimizedResponseGameID && minimizedResponseGameID == currentGameID)
        e.detail.headers["Minimal-Response"] = ""
});

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

window.addEventListener("htmx:wsAfterMessage", (e: CustomEvent) => {
    try {
        var message = JSON.parse(e.detail.message)
    } catch (err) {
        return
    }

    switch (message.event) {
        case "ERROR":
            alert(`${message.payload.type}: ${message.payload.detail}`)
        case "DELETED_PDF":
            window.dispatchEvent(new CustomEvent("deleted-pdf", { detail: { pdfID: message.payload.id } }))
            window.dispatchEvent(new CustomEvent("remove-tab", { detail: { tabID: message.payload.id } }))
            break;
    }
});
