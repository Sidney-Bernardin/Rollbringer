window.Alpine.data("tabContainer", () => ({
  currentTab: "",
  closeTab(name: string) {
    const tabButtons = (this.$el as HTMLElement).querySelectorAll(
      `.tab-button[data-tab-name=${name}]`,
    );
    tabButtons.forEach((elem) => elem.remove());
  },
}));
