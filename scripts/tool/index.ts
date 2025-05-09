import dotenv from "dotenv";
import { Logger } from "@m234/logger";

dotenv.config();

export const logClientComp = new Logger("client-compilation");
export const logInitFiles = new Logger("init-files");

export const encoder = new TextEncoder();
export const decoder = new TextDecoder("utf-8");

export enum envKeys {
    PORT = "PORT",
    IMAGE_CACHE_DURATION = "IMAGE_CACHE_DURATION",
}
