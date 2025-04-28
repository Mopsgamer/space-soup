import * as esbuild from "esbuild";
import { copy as copyPlugin } from "esbuild-plugin-copy";
import { existsSync, FSWatcher, watch } from "node:fs";
import { logClientComp } from "./tool/index.ts";
import tailwindcssPlugin from "esbuild-plugin-tailwindcss";
import { dirname } from "node:path";
import process from "node:process";

const folder = "client";
const args = process.argv.slice(2)
const isWatch = args.includes("watch");

type BuildOptions = esbuild.BuildOptions & {
    whenChange?: string[];
};

const minify = args.includes("min");

const options: esbuild.BuildOptions = {
    bundle: true,
    minify: minify,
    minifyIdentifiers: minify,
    minifySyntax: minify,
    minifyWhitespace: minify,
    platform: "browser",
    format: "esm",
    target: [
        "esnext",
        "chrome67",
        "edge79",
        "firefox68",
        "safari14",
    ],
};

let buildCalls = 0;
async function build(
    options: BuildOptions,
): Promise<void> {
    const { outdir, outfile, entryPoints = [], whenChange = [] } = options;
    buildCalls++;

    const directory = outdir || dirname(outfile!);
    if (calls.length == 1) {
        logClientComp.start("Bundling: %s", directory);
    } else {
        logClientComp.start(
            "Bundling %d/%d: %s",
            buildCalls,
            calls.length,
            directory,
        );
    }

    const entryPointsNormalized = Array.isArray(entryPoints)
        ? entryPoints
        : Object.keys(entryPoints);

    const badEntryPoints = entryPointsNormalized.filter((entry) => {
        const pth = typeof entry === "string" ? entry : entry.in;

        if (pth.includes("*")) {
            return false;
        }

        try {
            return !existsSync(pth);
        } catch {
            return true;
        }
    });

    if (badEntryPoints.length > 0) {
        logClientComp.error(
            `File expected to exist: ${badEntryPoints.join(", ")}`,
        );
        return;
    }

    if (!outfile && !outdir) {
        logClientComp.error(`Provide outdir or outfile.`);
        return;
    }

    if (outfile && outdir) {
        logClientComp.error(`Can not use outdir and outfile at the same time.`);
        return;
    }

    const safeOptions = options;
    delete safeOptions.whenChange;
    const ctx = await esbuild.context(safeOptions as esbuild.BuildOptions);

    async function rebuild() {
        try {
            await ctx.rebuild();
        } catch (error) {
            logClientComp.error(error);
        }
    }

    await rebuild();
    logClientComp.end(true);

    if (!isWatch) {
        await ctx.dispose();
        return;
    }

    try {
        await ctx.watch();
    } catch (error) {
        logClientComp.error(error);
        return;
    }

    if (whenChange.length === 0) {
        logClientComp.error("Nothing to watch: " + whenChange.join(", ") + ".");
        await ctx.dispose();
        return;
    }

    for (const p of whenChange) {
        let watcher: FSWatcher;
        try {
            watcher = watch(p, { recursive: true });
        } catch (error) {
            logClientComp.error(error);
            logClientComp.error("Bad paths, can not add watcher: " + p + ".");
            continue;
        }

        watcher.addListener("change", (kind) => {
            if (
                !(
                    kind === "modify" || kind === "create" ||
                    kind === "remove"
                )
            ) return;

            rebuild();
        });
        watcher.addListener("close", (kind) => {
            ctx.dispose();
        });
    }
}

function copy(from: string, to: string): Promise<void> {
    return build({
        ...options,
        outdir: to,
        entryPoints: [],
        whenChange: [to],
        plugins: [copyPlugin({
            once: isWatch,
            resolveFrom: "cwd",
            assets: { to: to + "/**/*", from },
            copyOnStart: true,
        })],
    });
}

// deno-lint-ignore no-explicit-any
type Call<Args extends (...args: any[]) => Promise<void>> = [
    fn: (...args: Parameters<Args>) => Promise<void>,
    params: Parameters<Args>,
    group: string[],
];

const slAlias = ["shoelace", "shoe", "sl"];

const calls: (Call<typeof copy> | Call<typeof build>)[] = [
    [copy, [
        "./node_modules/@shoelace-style/shoelace/dist/assets",
        `./${folder}/static/shoelace/assets`,
    ], [...slAlias]],

    [copy, [
        `./${folder}/src/assets`,
        `./${folder}/static/assets`,
    ], ["assets"]],

    [build, [{
        ...options,
        outdir: `./${folder}/static/js`,
        entryPoints: [`./${folder}/src/ts/**/*`],
        whenChange: [
            `./${folder}/static/js`,
        ],
    }], ["js", ...slAlias]],

    [build, [{
        ...options,
        outdir: `./${folder}/static/css`,
        entryPoints: [`./${folder}/src/tailwindcss/**/*.css`],
        whenChange: [
            `./${folder}/templates`,
            `./${folder}/src/tailwindcss`,
        ],
        external: ["/static/assets/*"],
        plugins: [
            tailwindcssPlugin(),
        ],
    }], ["css", ...slAlias]],
];

const existingGroups = Array.from(new Set(calls.flatMap((c) => c[2])));
const extraGroups = ["min", "watch", "all", "help"];
const availableGroups = [...extraGroups, ...existingGroups];

if (args.includes("help")) {
    logClientComp.info(
        "Available options: %s.",
        availableGroups.join(", "),
    );
    logClientComp.info(
        "Usage example:\n\n\tdeno task compile:client js css min watch\n",
    );
    process.exit();
}

const unknownGroups = args.filter(
    (a) => !availableGroups.includes(a),
);
if (unknownGroups.length > 0) {
    logClientComp.warn(
        `Unknown groups: ${unknownGroups.join(", ")}\n` +
            "Available groups: %s.",
        availableGroups.join(", "),
    );
}

logClientComp.info(
    `Starting bundling "./${folder}" ${isWatch ? " in watch mode" : ""}...`,
);

const existingGroupsUsed = !args.includes("all") &&
    existingGroups.some((g) => args.includes(g));

if (existingGroupsUsed) {
    calls.splice(
        0,
        calls.length,
        ...calls.filter(
            ([, , groups]) => {
                return groups.some((g) => {
                    const includes = args.includes(g);
                    return includes;
                });
            },
        ),
    );
}

for (const [fn, args] of calls) {
    // deno-lint-ignore no-explicit-any
    await fn(...args as any);
}

if (logClientComp.someFailed) {
    logClientComp.error("Bundled");
} else {
    logClientComp.success("Bundled successfully");
}
if (isWatch) {
    if (logClientComp.someFailed) {
        logClientComp.error("Watching for file changes...");
    } else {
        logClientComp.success("Watching for file changes...");
    }
} else if (logClientComp.someFailed) {
    process.exit(1);
}
