import Alpine from "alpinejs";

Alpine.data("modal", () => ({
    close(e: MouseEvent) {
        const dimentions: DOMRect = this.$el.getBoundingClientRect();

        if (
            e.clientX < dimentions.left ||
            e.clientX > dimentions.right ||
            e.clientY < dimentions.top ||
            e.clientY > dimentions.bottom
        ) this.$root.close()
    },
}))
