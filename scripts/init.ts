import dotenv from "dotenv";
import { existsSync } from "@std/fs";
import { decoder, encoder, envKeys, logInitFiles } from "./tool/index.ts";

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
    defaultEnv.set(envKeys.IMAGE_CACHE_DURATION, {
        value: "10m",
        comment:
            "images cache expiration time.\nfallback for invalid values - 10m.\nmemory is freed only when new images created.\nsee https://pkg.go.dev/time#ParseDuration",
    });

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
}

try {
    const path = ".env";
    logInitFiles.start(`Initializing '${path}'`);
    initEnvFile(path);
    logInitFiles.end(true);
} catch (error) {
    logInitFiles.error(error);
    Deno.exit(1);
}
