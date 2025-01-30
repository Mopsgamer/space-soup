import consola from "consola";

consola.options.formatOptions.columns = 3;

export const logBuild = consola.withTag("build");
export const logCleanup = consola.withTag("cleanup");
export const logInitDb = consola.withTag("init-database");
export const logInitFiles = consola.withTag("init-files");

export const encoder = new TextEncoder();
export const decoder = new TextDecoder("utf-8");

export enum envKeys {
    /**
     * Can be 0 (test), 1 (dev) or 2 (prod)
     * @default 1
     */
    ENVIRONMENT = "ENVIRONMENT",
    PORT = "PORT",
}
