import { Redis } from 'ioredis'

export * from './utils/user-rate-limiter'

export * from './utils/middleware'

export * from './fixed-window'
export * from './token-bucket'
// Redis client creation utility
export interface RedisConfig {
  host?: string
  port?: number
  url?: string
}

export function createRedisClient({
  host,
  port,
  user,
  password,
}: {
  host: string
  port: number
  user: string
  password: string
}): Redis {
  return new Redis(port, host, { username: user, password })
}
