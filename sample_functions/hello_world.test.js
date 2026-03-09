import { test } from "node:test";
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const __dirname = dirname(fileURLToPath(import.meta.url));

test("hello_world.js outputs 'Hello World'", (_, done) => {
  const filePath = join(__dirname, "hello_world.js");
  execFile("node", [filePath], (error, stdout, stderr) => {
    assert.ifError(error);
    assert.equal(stdout.trim(), "Hello World");
    done();
  });
});
