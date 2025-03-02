import { log } from '@repo/logger'

import { createServer } from './server'

const port = parseInt(process.env.PORT || '3003', 10)
const server = createServer()
log(`api running on ${port}`)

export default {
  port,
  fetch: server.fetch,
}
