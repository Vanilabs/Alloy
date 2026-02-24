import { makeLogger } from "./logger";

export const dlog = process.env.NODE_ENV === 'production' ? ((..._args: any[]) => {}):makeLogger('DEV');
export const plog = makeLogger('PROD');
export const log = makeLogger('DEFAULT');
