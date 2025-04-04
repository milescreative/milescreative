import { z } from 'zod'

import { redis } from './utils/redis'
import { rtxn_exec } from './utils/redis-helpers'

export class LeakyBucketRateLimiter {
  private bucketCapacity: number
  private leakRate: number
  private prefix: string

  constructor({
    bucketCapacity,
    leakRate,
    prefix,
  }: {
    bucketCapacity: number
    leakRate: number
    prefix: string
  }) {
    this.bucketCapacity = bucketCapacity
    this.leakRate = leakRate
    this.prefix = prefix
  }

  async isAllowed(key: string): Promise<boolean> {
    const prefix = this.prefix ? `${this.prefix}:` : ''
    const keyCount = `${prefix}${key}:lb_count`
    const keyLastLeak = `${prefix}${key}:lb_lastLeak`

    const currentTime = Date.now()

    // Current State
    const transaction = redis.multi()
    transaction.get(keyLastLeak)
    transaction.get(keyCount)
    const [lastLeakTime, lastRequestCount] = await rtxn_exec(transaction, [
      z.coerce.number().default(currentTime),
      z.coerce.number().default(0),
    ])
    let requestCount = lastRequestCount

    const elapsedTimeMs = currentTime - lastLeakTime
    const elaspedTimeSecs = elapsedTimeMs / 1000
    const requestsToLeak = Math.floor(elaspedTimeSecs * this.leakRate)
    requestCount = Math.max(requestCount - requestsToLeak, 0)

    const isAllowed = requestCount < this.bucketCapacity

    if (isAllowed) {
      requestCount++
    }

    // Update State
    transaction.set(keyLastLeak, currentTime.toString())
    transaction.set(keyCount, requestCount.toString())
    await transaction.exec()

    return isAllowed
  }
}
