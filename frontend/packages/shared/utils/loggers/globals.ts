import { dlog, plog, log } from './logger.runtime';

declare global {
    // Make them callable anywhere
    // eslint-disable-next-line no-var
    var $dlog: typeof dlog
    // eslint-disable-next-line no-var
    var $plog: typeof plog
    // eslint-disable-next-line no-var
    var $log: typeof log
}

// Attach once (safe for hot reload)
if (!globalThis.$dlog) globalThis.$dlog = dlog
if (!globalThis.$plog) globalThis.$plog = plog
if (!globalThis.$log) globalThis.$log = log

export { }
