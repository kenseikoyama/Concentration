import "./wasm_exec.js";

const path = new URL("main.wasm", import.meta.url);

const go = new Go();
const { instance } = await WebAssembly.instantiateStreaming(
  fetch(path),
  go.importObject,
);
go.run(instance);