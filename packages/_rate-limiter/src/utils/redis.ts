import { Redis } from 'ioredis'

const getRedis = () => {
  if (!process.env.REDIS_HOST) {
    throw new Error('REDIS_HOST is not set')
  }
  if (!process.env.REDIS_PORT) {
    throw new Error('REDIS_PORT is not set')
  }
  if (!process.env.REDIS_USER) {
    throw new Error('REDIS_USER is not set')
  }
  if (!process.env.REDIS_PASSWORD) {
    throw new Error('REDIS_PASSWORD is not set')
  }

  const r = new Redis(
    parseInt(process.env.REDIS_PORT || '4545'),
    process.env.REDIS_HOST || '192.168.1.50',
    {
      username: process.env.REDIS_USER || 'default',
      password: process.env.REDIS_PASSWORD || '',
    },
  )

  r.on('connect', () => {
    console.log('Connected to Redis')
  })

  r.on('error', (err) => {
    console.error('Redis connection error:', err)
  })

  r.on('connecting', () => {
    console.log('Connecting to Redis')
  })

  r.on('ready', () => {
    console.log('Redis connection ready')
  })

  r.on('wait', () => {
    console.log('Redis connection waiting')
  })

  return r
}

export const redis = getRedis()
