import Alpine from "alpinejs";

Alpine.bind("tabButton", (tabName: string) => ({
  "x-bind:class"() {
    return this.$data.currentTab == tabName && "active";
  },

  "x-on:click"() {
    this.$data.currentTab = tabName;
  },
}));
