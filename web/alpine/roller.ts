window.alpine.data("roller", () => ({
    dice: [],

    addDie(e: Event) {
        const target = e.target as HTMLSelectElement;

        if (target.value !== "default") {
            this.dice.push(target.value);
            target.value = "default";
        }
    },

    removeDie(die: string) {
        this.dice.splice(this.dice.indexOf(die), 1);
    },
}));
