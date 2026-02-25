export type LogChannel = 'DEV' | 'PROD' | 'DEFAULT' | 'ALL';
export type LogLevel = 'INFO' | 'WARN' | 'ERROR';
export type ChannelRule = { channel: LogChannel; namespace?: string };
export type AllRule = { channel: 'ALL' };

export type AllowedRule = ChannelRule | AllRule;

type AnyArgs = any[]

const LEVEL_RANK: Record<LogLevel, number> = {
    INFO: 1,
    WARN: 2,
    ERROR: 3
}

function parseLogLevel(raw: string | undefined): LogLevel {
    const v = (raw || '').trim().toUpperCase() as LogLevel;
    if (v === 'WARN') return 'WARN';
    if (v === 'ERROR') return 'ERROR';
    return 'INFO';
}

function isChannelRule(r: AllowedRule): r is ChannelRule {
    return r.channel !== 'ALL'
}

function parseAllowedLogs(raw: string | undefined): AllowedRule[] {
    if (!raw) return []
    const parts = raw
        .split(',')
        .map(s => s.trim())
        .filter(Boolean)

    const rules: AllowedRule[] = []
    for (const p of parts) {
        const upper = p.toUpperCase()
        if (upper === 'ALL') {
            rules.push({ channel: 'ALL' })
            continue
        }

        // Examples: DEV, DEV:auth, PROD:payments, DEFAULT:boot
        const [chRaw, nsRaw] = p.split(':')
        const channel = chRaw!.trim().toUpperCase() as LogChannel
        const namespace = nsRaw?.trim()

        if (channel === 'DEV' || channel === 'PROD' || channel === 'DEFAULT') {
            rules.push({ channel, namespace })
        }
    }

    return rules
}

const allowedRules = parseAllowedLogs(process.env.NEXT_PUBLIC_ALLOWED_LOGS)
const minLevel = parseLogLevel(process.env.NEXT_PUBLIC_LOG_LEVEL)

function isAllowed(channel: LogChannel, namespace?: string) {
    // ALL means everything
    if (allowedRules.some(r => r.channel === 'ALL')) return true

    // Match channel-only: DEV
    if (allowedRules.some(r => isChannelRule(r) && r.channel === channel && !r.namespace)) return true

    // Match channel+namespace: DEV:auth
    if (namespace) {
        // return allowedRules.some(r=>r.channel === channel && 'namespace' in r && r.namespace === namespace);
        return allowedRules.some(r => isChannelRule(r) && r.channel === channel && r.namespace === namespace)
    }

    return false
}

function levelAllowed(level: LogLevel) {
    return LEVEL_RANK[level] >= LEVEL_RANK[minLevel]
}

function formatPrefix(channel: LogChannel, level: LogLevel, namespace?: string) {
    const ns = namespace ? `:${namespace}` : ''
    return `[${channel}${ns}] [${level}]`
}

function emit(level: LogLevel, prefix: string, args: AnyArgs) {
    if (level === 'ERROR') console.error(prefix, ...args)
    else if (level === 'WARN') console.warn(prefix, ...args)
    else console.log(prefix, ...args)
}

/**
 * Core factory.
 * Usage:
 *   dlog('auth', 'token refreshed')
 *   plog('payments', 'charge ok')
 *   log('boot', 'runtime init')
 */
export function makeLogger(channel: LogChannel) {
    return (namespaceOrMsg: string, ...rest: AnyArgs) => {
        // If user passes only one arg like dlog("hello"), treat as msg without namespace
        let namespace: string | undefined
        let msgArgs: AnyArgs

        if (rest.length === 0) {
            namespace = undefined
            msgArgs = [namespaceOrMsg]
        } else {
            namespace = namespaceOrMsg
            msgArgs = rest
        }

        // If first of msgArgs is a level, allow it:
        // dlog('auth', 'WARN', 'something')
        // dlog('WARN', 'something')  -> no namespace
        let level: LogLevel = 'INFO'
        if (typeof msgArgs[0] === 'string') {
            const maybe = msgArgs[0].toUpperCase()
            if (maybe === 'INFO' || maybe === 'WARN' || maybe === 'ERROR') {
                level = maybe as LogLevel
                msgArgs = msgArgs.slice(1)
            }
        }

        if (!levelAllowed(level)) return
        if (!isAllowed(channel, namespace)) return

        emit(level, formatPrefix(channel, level, namespace), msgArgs)
    }
}
