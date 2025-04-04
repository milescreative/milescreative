import { test_result } from '@milescreative/db'
import { log } from '@milescreative/logger'
import {
  FixedWindowRateLimiter,
  getClientIp,
  LeakyBucketRateLimiter,
  TokenBucketRateLimiter,
} from '@milescreative/rate-limiter'
import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'

export const createServer = (): Hono => {
  const app = new Hono()

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
  const leakyBucketRateLimiter = new LeakyBucketRateLimiter({
    bucketCapacity: 5,
    leakRate: 1,
    prefix: 'myapp',
  })

  app
    .use('*', logger())
    .use('*', cors())
    .use('*', async (c, next) => {
      const ip =
        process.env.NODE_ENV !== 'production'
          ? '127.0.0.1'
          : getClientIp(c.req.header())
      if (!ip) {
        return c.json({ message: 'No IP address found' }, 400)
      }
      const isAllowed = await tokenBucketRateLimiter.isAllowed(ip)
      if (!isAllowed) {
        return c.json(
          { message: 'User rate limit exceeded. Please try again later.' },
          429
        )
      }
      const global_isAllowed = await leakyBucketRateLimiter.isAllowed('')
      if (!global_isAllowed) {
        return c.json(
          { message: 'Experiencing high load, please try again later' },
          429
        )
      }
      return next()
    })
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
    .get('/rate-limit-leaky-bucket', async (c) => {
      const isAllowed = await leakyBucketRateLimiter.isAllowed('test')
      if (!isAllowed) {
        return c.json({ message: 'Rate limit exceeded' }, 429)
      }
      return c.json({ message: 'This page is rate limited' })
    })
  return app
}
