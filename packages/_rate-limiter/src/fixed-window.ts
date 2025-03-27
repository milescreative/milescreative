import { z } from 'zod'

import { redis } from './redis'

export class FixedWindowRateLimiter {
  private limit: number
  private windowSize: string
  private prefix: string

  constructor({
    limit,
    windowSize,
    prefix,
  }: {
    limit: number
    windowSize: string
    prefix: string
  }) {
    this.limit = limit
    this.windowSize = windowSize
    this.prefix = prefix
  }

  async isAllowed(key: string): Promise<boolean> {
    const prefix = this.prefix ? `${this.prefix}:` : ''
    const currentKey = `${prefix}${key}`
    const currentCount = z.coerce
      .number()
      .default(0)
      .parse(await redis.get(currentKey))
    const isAllowed = currentCount < this.limit

    if (isAllowed) {
      const transaction = redis.multi()
      transaction.incr(currentKey)
      transaction.expire(currentKey, this.windowSize, 'NX')
      await transaction.exec()
    }

    return isAllowed
  }
}
