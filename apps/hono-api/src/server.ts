import { test_result } from '@milescreative/db'
import { log } from '@milescreative/logger'
import {
  createRedisClient,
  createUserRateLimiter,
  FixedWindowRateLimiter,
  rateLimiterMiddleware,
  TokenBucketRateLimiter,
} from '@milescreative/rate-limiter'
import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'

export const createServer = (): Hono => {
  const app = new Hono()
  // Create Redis client using the factory from the rate-limiter package
  const redis = createRedisClient({
    host: process.env.REDIS_HOST || 'localhost',
    port: parseInt(process.env.REDIS_PORT || '6379', 10),
    user: process.env.REDIS_USER || 'default',
    password: process.env.REDIS_PASSWORD || '',
  })

  const rateLimiter = createUserRateLimiter(
    {
      redis,
      appPrefix: 'myapp',
    },
    {
      maxRequests: 10, // Allow 100 requests
      windowSeconds: 10, // Per minute
    }
  )
  const fixedWindowRateLimiter = new FixedWindowRateLimiter({
    limit: 10,
    windowSize: '10',
    prefix: 'myapp',
  })
  const tokenBucketRateLimiter = new TokenBucketRateLimiter({
    bucketCapacity: 5,
    refillRate: 1,
    prefix: 'myapp',
  })

  app
    .use('*', logger())
    .use('*', cors())
    // .use('/rate-limit', async (c, next) => {
    //   const headers = c.req.header()
    //   let modifiedHeaders = headers
    //   if (process.env.NODE_ENV !== 'production') {
    //     modifiedHeaders = {
    //       ...headers,
    //       'x-forwarded-for': '127.0.0.1',
    //     }
    //   }

    //   log(`Rate limit middleware`)
    //   log(JSON.stringify(modifiedHeaders))
    //   return rateLimiterMiddleware(
    //     { headers: modifiedHeaders, res: c.res },
    //     next,
    //     rateLimiter
    //   )
    // })
    .get('/message/:name', (c) => {
      const name = c.req.param('name')
      return c.json({ message: `hello ${name}` })
    })
    .get('/status', (c) => {
      return c.json({ ok: true })
    })
    .get('/test', async (c) => {
      const result = await test_result()
      return c.json(result)
    })
    .get('/env', async (c) => {
      const result = JSON.stringify(process.env)
      return c.json(result)
    })
    .get('/rate-limit', async (c) => {
      const isAllowed = await fixedWindowRateLimiter.isAllowed('test')
      if (!isAllowed) {
        return c.json({ message: 'Rate limit exceeded' }, 429)
      }
      return c.json({ message: 'This page is rate limited' })
    })
    .get('/rate-limit-token-bucket', async (c) => {
      const isAllowed = await tokenBucketRateLimiter.isAllowed('test')
      if (!isAllowed) {
        return c.json({ message: 'Rate limit exceeded' }, 429)
      }
      return c.json({ message: 'This page is rate limited' })
    })
  return app
}
