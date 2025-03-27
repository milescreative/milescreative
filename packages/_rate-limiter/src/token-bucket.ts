import { z } from 'zod'

import { redis } from './redis'
import { rtxn_exec } from './utils/redis-helpers'

export class TokenBucketRateLimiter {
  private bucketCapacity: number
  private refillRate: number
  private prefix: string

  constructor({
    bucketCapacity,
    refillRate,
    prefix,
  }: {
    bucketCapacity: number
    refillRate: number
    prefix: string
  }) {
    this.bucketCapacity = bucketCapacity
    this.refillRate = refillRate
    this.prefix = prefix
  }

  async isAllowed(key: string): Promise<boolean> {
    const prefix = this.prefix ? `${this.prefix}:` : ''
    const keyCount = `${prefix}${key}:count`
    const keyLastRefill = `${prefix}${key}:lastRefill`

    const currentTime = Date.now()

    // Current State
    const transaction = redis.multi()
    transaction.get(keyLastRefill)
    transaction.get(keyCount)
    const [lastRefillTime, lasttokenCount] = await rtxn_exec(transaction, [
      z.coerce.number().default(currentTime),
      z.coerce.number().default(this.bucketCapacity),
    ])
    let tokenCount = lasttokenCount

    const elapsedTimeMs = currentTime - lastRefillTime
    const elaspedTimeSecs = elapsedTimeMs / 1000
    const tokensToAdd = Math.floor(elaspedTimeSecs * this.refillRate)
    tokenCount = Math.min(tokenCount + tokensToAdd, this.bucketCapacity)

    const isAllowed = tokenCount >= 0

    if (isAllowed) {
      tokenCount--
    }

    // Update State
    transaction.set(keyLastRefill, currentTime.toString())
    transaction.set(keyCount, tokenCount.toString())
    await transaction.exec()

    return isAllowed
  }
}
