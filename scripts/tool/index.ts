import dotenv from "dotenv";
import { Logger } from "./logger.ts";

dotenv.config();

export const logClientComp = new Logger("client-compilation");
export const logInitFiles = new Logger("init-files");

export const encoder = new TextEncoder();
export const decoder = new TextDecoder("utf-8");

export enum envKeys {
    PORT = "PORT",
}
