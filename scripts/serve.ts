import { logClientComp } from "./tool/index.ts";
import kill from "tree-kill";

const paths = ["server", "main.go"];

const serverCommand = new Deno.Command("go", {
    args: ["run", "."],
});
let goRunProcess: Deno.ChildProcess | undefined = undefined;

function start() {
    goRunProcess = serverCommand.spawn();
}

async function watchAndRestart() {
    start();
    const watcher = Deno.watchFs(paths, { recursive: true });
    let prevKind = ""
    let prevTime = 0
    for await (const event of watcher) {
        if (
            !(
                event.kind === "modify" || event.kind === "create" ||
                event.kind === "remove"
            )
        ) continue;
        const nowTime = Date.now()
        if (event.kind == prevKind && (nowTime - prevTime < 300)) continue;

        tryToKill();
        logClientComp.info(
            "File change detected: %s. Restarting...",
            event.kind,
        );
        prevKind = event.kind
        prevTime = nowTime
        start();
    }
}

function tryToKill() {
    if (goRunProcess == undefined) {
        return;
    }
    try {
        kill(goRunProcess.pid, "SIGTERM");
    } catch { /* empty */ }
    goRunProcess = undefined;
}

watchAndRestart();
