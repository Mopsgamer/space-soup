import dotenv from "dotenv";
import { Logger } from "./logger.ts";

dotenv.config();

export const logClientComp = new Logger("client-compilation");
export const logInitDb = new Logger("init-database");
export const logInitFiles = new Logger("init-files");

export const encoder = new TextEncoder();
export const decoder = new TextDecoder("utf-8");

export enum envKeys {
    JWT_KEY = "JWT_KEY",
    USER_AUTH_TOKEN_EXPIRATION = "USER_AUTH_TOKEN_EXPIRATION",
    CHAT_MESSAGE_MAX_LENGTH = "CHAT_MESSAGE_MAX_LENGTH",

    PORT = "PORT",

    DB_PASSWORD = "DB_PASSWORD",
    DB_NAME = "DB_NAME",
    DB_USER = "DB_USER",
    DB_HOST = "DB_HOST",
    DB_PORT = "DB_PORT",
}
