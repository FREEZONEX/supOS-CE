declare module 'js-yaml' {
  export function load(input: string, options?: unknown): unknown;
  export function dump(input: unknown, options?: unknown): string;

  const jsYaml: {
    load: typeof load;
    dump: typeof dump;
  };

  export default jsYaml;
}
