import { watch } from "node:fs";
import { type ChildProcess, spawn } from "node:child_process";
import { logClientComp } from "./tool/index.ts";
import kill from "tree-kill";

const paths = ["server", "main.go"];

let goRunProcess: ChildProcess | undefined;

function start() {
    goRunProcess = spawn("go", ["run", "lite.go", "-tags=lite"], {stdio: "inherit"});
}

async function watchAndRestart() {
    start();
    for (const p of paths) {
        const watcher = watch(p, { recursive: true });
        let prevKind = "";
        let prevTime = 0;
        watcher.addListener("change", (kind) => {
            if (
                !(
                    kind === "modify" || kind === "create" ||
                    kind === "remove"
                )
            ) return;
            const nowTime = Date.now();
            if (kind == prevKind && (nowTime - prevTime < 300)) return;

            tryToKill();
            logClientComp.info(
                "File change detected: %s. Restarting...",
                kind,
            );
            prevKind = kind;
            prevTime = nowTime;
            start();
        });
    }
}

function tryToKill() {
    if (goRunProcess?.pid == undefined) {
        return;
    }
    try {
        kill(goRunProcess.pid, "SIGTERM");
    } catch { /* empty */ }
    goRunProcess = undefined;
}

watchAndRestart();
