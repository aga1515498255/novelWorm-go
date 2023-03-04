export let config = {
  mode: "dev",
  devBackendPrefix: "http://127.0.0.1:4321",

  getPrefix() {
    if (this.mode === "dev") {
      return this.devBackendPrefix;
    } else {
      return "";
    }
  },
};
