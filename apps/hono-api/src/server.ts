import { test_result } from '@repo/db'
import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'

export const createServer = (): Hono => {
  const app = new Hono()

  app
    .use('*', logger())
    .use('*', cors())
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

  return app
}
