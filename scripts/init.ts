import dotenv from "dotenv";
import { decoder, encoder, envKeys, logInitFiles } from "./tool/index.ts";
import { existsSync, readFileSync, writeFileSync } from "node:fs";
import process from "node:process";

function initEnvFile(path: string): void {
    type EnvKeyEntry = {
        value?: string | number | boolean;
        comment?: string;
    };
    const defaultEnv = new Map<string, EnvKeyEntry>();

    defaultEnv.set(envKeys.PORT, {
        value: 3000,
        comment: "application port",
    });

    const env = existsSync(path)
        ? dotenv.parse(decoder.decode(readFileSync(path)))
        : {};

    writeFileSync(
        path,
        encoder.encode(
            Array.from(defaultEnv.entries()).map(
                ([key, { value, comment }]) => {
                    env[key] ||= value === undefined ? "" : String(value);
                    process.env[key] = env[key];
                    if (value == undefined) {
                        comment += "\ndefault: <empty>";
                    } else {
                        comment += "\ndefault: " + value;
                    }
                    return `${
                        comment
                            ? "# " + comment.replaceAll("\n", "\n# ") + "\n"
                            : ""
                    }${key}=${env[key]}\n\n`;
                },
            ).join(""),
        ),
    );
}

try {
    const path = ".env";
    logInitFiles.start(`Initializing '${path}'`);
    initEnvFile(path);
    logInitFiles.end(true);
} catch (error) {
    logInitFiles.error(error);
    process.exit(1);
}
