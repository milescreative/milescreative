import { log } from '@milescreative/logger'

import { getClientIp } from '../ip'
import { UserRateLimiter } from '../user-rate-limiter'

export type Next = () => void | Promise<void>

export type Context = {
  headers: Request['headers'] | Record<string, string>
  res: Response
}

export async function rateLimiterMiddleware(
  context: Context,
  next: Next,
  rateLimiter: UserRateLimiter,
  debug: boolean = false
) {
  const userId = getClientIp(context.headers)
  if (!userId) {
    log(`Context: ${JSON.stringify(context)}`)
    return new Response('Unauthorized - No IP address found', { status: 401 })
  }
  const result = await rateLimiter.isAllowed(userId)
  if (!result.allowed) {
    return new Response('Rate limit exceeded', { status: 429 })
  }
  if (debug) {
    return new Response(result.remaining.toString(), { status: 200 })
  }
  return next()
}
