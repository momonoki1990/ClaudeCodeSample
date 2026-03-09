describe("hello_world.js", () => {
  beforeEach(() => {
    jest.resetModules();
    jest.spyOn(console, "log").mockImplementation(() => {});
  });

  afterEach(() => {
    console.log.mockRestore();
  });

  test("'Hello World' と出力する", () => {
    require("./hello_world");
    expect(console.log).toHaveBeenCalledWith("Hello World");
  });
});
