import { Redis } from 'ioredis'

export const redis = new Redis(
  parseInt(process.env.REDIS_PORT || '4545'),
  process.env.REDIS_HOST || 'localhost',
  {
    username: process.env.REDIS_USER || 'default',
    password: process.env.REDIS_PASSWORD || '',
  }
)
console.log(redis.status)

// Add connection event handlers
redis.on('connect', () => {
  console.log('Connected to Redis')
})

redis.on('error', (err) => {
  console.error('Redis connection error:', err)
})
