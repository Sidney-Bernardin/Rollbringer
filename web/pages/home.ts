function button() {
    const elem: HTMLDivElement = document.createElement('div')

    elem.innerHTML = "Hello, World!"

    return elem
}

document.body.appendChild(button())
