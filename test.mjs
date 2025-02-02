import {readFileSync} from 'node:fs'
import vm from 'vm';

vm.runInThisContext(readFileSync("vendor/wasm_exec.js", 'utf-8'));

const go = new Go();

const result = await WebAssembly.instantiate(readFileSync('test.wasm'), go.importObject)
await go.run(result.instance)
