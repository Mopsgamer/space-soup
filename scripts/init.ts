import dotenv from "dotenv";
// @deno-types="npm:@types/mysql"
import mysql from "mysql2";
import { existsSync } from "@std/fs";
import {
    decoder,
    encoder,
    envKeys,
    logInitFiles,
} from "./tool/index.ts";

function initEnvFile(): void {
    type EnvKeyEntry = {
        value?: string | number | boolean;
        comment?: string;
    };
    const defaultEnv = new Map<string, EnvKeyEntry>();

    defaultEnv.set(envKeys.PORT, {
        value: 3000,
        comment: "application port",
    });

    const path = ".env";
    const env = existsSync(path)
        ? dotenv.parse(decoder.decode(Deno.readFileSync(path)))
        : {};

    Deno.writeFileSync(
        path,
        encoder.encode(
            Array.from(defaultEnv.entries()).map(
                ([key, { value, comment }]) => {
                    env[key] ||= value === undefined ? "" : String(value);
                    Deno.env.set(key, env[key]);
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

    logInitFiles.success("Writed " + path);
}

try {
    initEnvFile();
} catch (error) {
    logInitFiles.error(error);
    Deno.exit(1);
}
