import { describe, test, expect, vi, beforeEach, afterEach } from "vitest";

describe("hello_world.js", () => {
  beforeEach(() => {
    vi.spyOn(console, "log").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  test("'Hello World' と出力する", async () => {
    await import("./hello_world.js");
    expect(console.log).toHaveBeenCalledWith("Hello World");
  });
});
