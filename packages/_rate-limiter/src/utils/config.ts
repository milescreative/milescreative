import Redis from 'ioredis'

export interface RedisConfig {
  host: string
  port: number
}

export function createRedisClient(config: RedisConfig): Redis {
  return new Redis(config.port, config.host)
}
