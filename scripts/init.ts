import dotenv from "dotenv";
import { existsSync } from "@std/fs";
import { decoder, encoder, envKeys, logInitFiles } from "./tool.ts";

function initEnvFile(): void {
    type EnvKeyEntry = {
        value?: string;
        comment?: string;
    };
    const defaultEnv = new Map<string, EnvKeyEntry>();
    defaultEnv.set(envKeys.ENVIRONMENT, {
        value: "1",
        comment: "can be 0 (test), 1 (dev) or 2 (prod)\ndefault: 1",
    });
    defaultEnv.set(envKeys.PORT, {
        value: "3000",
        comment: "application port\ndefault: 3000",
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
                    env[key] ||= value ?? "";
                    Deno.env.set(key, env[key]);
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
    logInitFiles.fatal(error);
    Deno.exit(1);
}
