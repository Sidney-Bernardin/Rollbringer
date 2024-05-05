window.alpine.data("modal", () => ({
    tryClosing(e: MouseEvent) {
        const dimentions: DOMRect = this.$el.getBoundingClientRect();

        if (
            e.clientX < dimentions.left ||
            e.clientX > dimentions.right ||
            e.clientY < dimentions.top ||
            e.clientY > dimentions.bottom
        ) this.closeModal()
    },

    closeModal() {
        this.$root.close();
    },
}));
