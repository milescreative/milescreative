import chalk from 'chalk'

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

interface LoggerOptions {
  level: LogLevel
  prefix?: string
  timestamp?: boolean
  colorize?: boolean
}

const LOG_LEVELS: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
}

const createLogger = (initialOptions: Partial<LoggerOptions> = {}) => {
  const options: LoggerOptions = {
    level: 'debug',
    prefix: '',
    timestamp: true,
    colorize: true,
    ...initialOptions,
  }

  const shouldLog = (level: LogLevel): boolean => {
    return LOG_LEVELS[level] >= LOG_LEVELS[options.level]
  }

  const formatMessage = (level: LogLevel, message: string): string => {
    const parts: string[] = []

    if (options.timestamp) {
      parts.push(`[${new Date().toISOString()}]`)
    }

    if (options.prefix) {
      parts.push(`[${options.prefix}]`)
    }

    parts.push(`[${level.toUpperCase()}]`)
    parts.push(message)

    const output = parts.join(' ')

    if (options.colorize) {
      switch (level) {
        case 'debug':
          return chalk.cyan(output)
        case 'info':
          return chalk.green(output)
        case 'warn':
          return chalk.yellow(output)
        case 'error':
          return chalk.red(output)
        default:
          return output
      }
    }

    return output
  }

  const logger = (message: string, ...args: unknown[]): void => {
    if (shouldLog('debug')) {
      console.debug(formatMessage('debug', message), ...args)
    }
  }

  Object.assign(logger, {
    debug: (message: string, ...args: unknown[]): void => {
      if (shouldLog('debug')) {
        console.debug(formatMessage('debug', message), ...args)
      }
    },

    info: (message: string, ...args: unknown[]): void => {
      if (shouldLog('info')) {
        console.info(formatMessage('info', message), ...args)
      }
    },

    warn: (message: string, ...args: unknown[]): void => {
      if (shouldLog('warn')) {
        console.warn(formatMessage('warn', message), ...args)
      }
    },

    error: (message: string, ...args: unknown[]): void => {
      if (shouldLog('error')) {
        console.error(formatMessage('error', message), ...args)
      }
    },

    setLevel: (level: LogLevel): void => {
      options.level = level
    },

    setPrefix: (prefix: string): void => {
      options.prefix = prefix
    },

    setTimestamp: (enabled: boolean): void => {
      options.timestamp = enabled
    },

    setColorize: (enabled: boolean): void => {
      options.colorize = enabled
    },
  })

  return logger
}

export const log = createLogger()
export default createLogger
