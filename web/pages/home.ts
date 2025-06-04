import Alpine from 'alpinejs'

import "htmx.org"

/////

document.addEventListener("htmx:beforeSwap", (e: CustomEvent<{
    isError: boolean,
}>) => {
    if (e.detail.isError) alert(e.detail.serverResponse)
})

/////

Alpine.start()
