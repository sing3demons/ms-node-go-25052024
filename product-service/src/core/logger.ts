import { createLogger, format, transports, Logger as WinstonLog } from 'winston'
import Context from './context'

export type ContextType = Record<string, string>

export type ILogger = {
  info: (message: string, data?: {} | [], ctx?: ContextType) => void
  warn: (message: string, data?: {} | [], ctx?: ContextType) => void
  error: (message: string, data?: any, ctx?: ContextType) => void
  debug: (message: string, data?: {} | [], ctx?: ContextType) => void
}

interface ISensitive {
  maskNumber(mobileNo: string, mask?: string): string
  maskEmail(email: string): string
  maskPassword(password: string): string
  masking(item: any): void
}

interface IgnoreCase {
  equal(a: string, b: string): boolean
  notEqual(a: string, b: string): boolean
  contain(a: string, b: string): boolean
  notContain(a: string, b: string): boolean
  startWith(a: string, b: string): boolean
}

type ILog = {
  info: (message: string, extra?: object) => void
  error: (message: string, extra?: object) => void
  warn: (message: string, extra?: object) => void
  debug: (message: string, extra?: object) => void
}

function makeStructuredClone<T>(obj: T): T {
  if (typeof obj === 'undefined') {
    return obj
  }
  const payload = JSON.parse(JSON.stringify(obj))
  if (typeof payload === 'object') {
    if (Array.isArray(payload)) {
      for (const item of payload) {
        if (typeof item === 'object') {
          Sensitive.masking(item)
        }
      }
    } else {
      Sensitive.masking(payload)
    }
  }
  return payload
}

const Sensitive: ISensitive = {
  maskNumber: (mobileNo: string, mask?: string): string => {
    let maskData = 'XXX-XXX-XX'
    if (mask) {
      maskData = maskData.replace(/X/g, mask)
    }
    if (ignoreCase.startWith(mobileNo, '+')) {
      if (mobileNo.length >= 10) {
        return `${maskData}${mobileNo.substring(mobileNo.length - 2, mobileNo.length)}`
      }
    } else if (ignoreCase.startWith(mobileNo, '0') && mobileNo.length >= 10) {
      return `${maskData}${mobileNo.substring(mobileNo.length - 2, mobileNo.length)}`
    }

    return mobileNo
  },
  maskEmail: (email: string): string => {
    const rex = new RegExp(/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/)
    if (!rex.test(email)) {
      return email
    } else {
      let [first, second] = email.split('@')
      if (!first) {
        return ''
      }
      if (first.length > 2) {
        const mask = first.substring(3, first.length)
        const notMask = first.substring(0, 3)
        first = notMask + 'X'.repeat(mask.length)
      } else {
        first = first.replace(first.substring(1, first.length), 'X'.repeat(first.length - 1))
      }
      return `${first}@${second}`
    }
  },
  maskPassword: (password: string): string => password.replace(password, '********'),
  masking: (item: any) => {
    for (const key in item) {
      if (ignoreCase.equal(key, 'password')) {
        item[key] = Sensitive.maskPassword(item[key])
      } else if (ignoreCase.equal(key, 'email')) {
        item[key] = Sensitive.maskEmail(item[key])
      } else if (ignoreCase.equal(key, 'mobileNo')) {
        item[key] = Sensitive.maskNumber(item[key])
      } else if (ignoreCase.equal(key, 'phone')) {
        item[key] = Sensitive.maskNumber(item[key])
      } else if (typeof item[key] === 'object') {
        Sensitive.masking(item[key])
      }
    }
  },
}

export const ignoreCase: IgnoreCase = {
  equal: (a?: string, b?: string) => {
    if (a === undefined || b === undefined) {
      return false
    }
    return a.toLowerCase() === b.toLowerCase()
  },
  notEqual: (a: string, b: string) => a.toLowerCase() !== b.toLowerCase(),
  contain: (a: string, b: string) => a.toLowerCase().includes(b.toLowerCase()),
  notContain: (a: string, b: string) => !a.toLowerCase().includes(b.toLowerCase()),
  startWith: (a: string, b: string) => a.toLowerCase().startsWith(b.toLowerCase()),
}

const level = process.env.LOG_LEVEL ?? 'debug'

class Logger implements ILogger {
  private readonly log: WinstonLog
  constructor() {
    this.log = createLogger({
      level: level,
      format: format.combine(
        format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss', alias: '@timestamp' }),
        format.json({
          replacer(key, value) {
            if (ignoreCase.equal(key, 'password')) {
              return Sensitive.maskPassword(value)
            } else if (ignoreCase.equal(key, 'email')) {
              return Sensitive.maskEmail(value)
            } else if (ignoreCase.equal(key, 'mobileNo')) {
              return Sensitive.maskNumber(value)
            } else if (ignoreCase.equal(key, 'phone')) {
              return Sensitive.maskPassword(value)
            } else if (key === 'timestamp') {
              return undefined
            }
            return value
          },
          // space: 2
        })
      ),
      exceptionHandlers: [],
      exitOnError: false,
      transports: [
        new transports.Console({
          level: level,
          handleExceptions: true,
        }),
      ],
      defaultMeta: { serviceName: process.env.SERVICE_NAME ?? 'Service-HTTP' },
    })
  }

  Logger(ctx: ContextType = Context.get(), extra?: object): ILogger {
    return this.log.child({ ...ctx, ...extra }) as ILogger
  }

  info(message: string, data?: {} | [], ctx?: ContextType) {
    const action = { ...makeStructuredClone(data) }
    this.log.info(message, { action, session: ctx?.session })
  }

  warn(message: string, data?: {} | [], ctx?: ContextType) {
    const action = { ...makeStructuredClone(data) }
    this.log.warn(message, { action, ...ctx })
  }

  error(message: string, data?: any, ctx?: ContextType) {
    const action = { ...makeStructuredClone(data) }
    this.log.error(message, { action, ...ctx })
  }

  debug(message: string, data?: {} | [], ctx?: ContextType) {
    const action = { ...makeStructuredClone(data) }
    this.log.debug(message, { action, ...ctx })
  }
}

export type LoggerType = Logger

export default Logger
