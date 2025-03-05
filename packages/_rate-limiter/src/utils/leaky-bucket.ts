import type { Redis } from 'ioredis'

interface RateLimiterConfig {
  redis: Redis
  appPrefix: string
}

interface LeakyBucketOptions {
  /**
   * Maximum capacity of the bucket (max requests that can be queued)
   */
  bucketCapacity: number

  /**
   * Leak rate in requests per second
   */
  leakRatePerSecond: number
}
interface LeakyBucketResult {
  /**
   * Whether the request is allowed based on bucket capacity
   */
  allowed: boolean

  /**
   * Current fill level of the bucket
   */
  currentLevel: number

  /**
   * Maximum capacity of the bucket
   */
  capacity: number

  /**
   * Estimated wait time in milliseconds if request is denied
   */
  estimatedWaitMs: number
}

/**
 * Leaky bucket rate limiter functions
 */
export interface LeakyBucketRateLimiter {
  /**
   * Check if a request can be processed based on the leaky bucket algorithm
   * @param bucketName Identifier for the specific bucket (e.g., API endpoint)
   * @returns Object containing bucket information and whether the request is allowed
   */
  isAllowed(bucketName: string): Promise<LeakyBucketResult>

  /**
   * Reset a specific bucket
   * @param bucketName Identifier for the specific bucket
   */
  reset(bucketName: string): Promise<void>
}

/**
 * Create a rate limiter using the leaky bucket algorithm
 *
 * @param config Configuration with Redis instance and app prefix
 * @param options Leaky bucket configuration options with defaults
 * @returns Functions to check rate limits and reset buckets
 */
export function createLeakyBucketRateLimiter(
  config: RateLimiterConfig,
  options: LeakyBucketOptions = {
    bucketCapacity: 1000,
    leakRatePerSecond: 10,
  }
): LeakyBucketRateLimiter {
  const redis = config.redis
  const prefix = `${config.appPrefix}:leaky-bucket:`

  /**
   * Check if a request can be processed based on the leaky bucket algorithm
   */
  async function isAllowed(bucketName: string): Promise<LeakyBucketResult> {
    const key = `${prefix}${bucketName}`
    const now = Date.now()

    // Use Redis transaction to ensure atomicity
    const multi = redis.multi()

    // Get current bucket state
    multi.hgetall(key)

    // Execute the pipeline to get current state
    const results = await multi.exec()
    if (!results) {
      throw new Error('Redis operation failed')
    }

    // Parse bucket state with safe defaults
    const bucketState = (results[0]?.[1] || {}) as Record<string, string>
    const lastUpdate = parseInt(bucketState.lastUpdate || `${now}`, 10)
    const currentLevel = parseInt(bucketState.level || '0', 10)

    // Calculate leakage since last update
    const timeElapsed = now - lastUpdate
    const leakRate = options.leakRatePerSecond
    const leaked = Math.floor((timeElapsed * leakRate) / 1000)
    const newLevel = Math.max(0, currentLevel - leaked)

    // Check if the bucket has capacity
    const allowed = newLevel < options.bucketCapacity

    // Calculate new level based on whether request is allowed
    const finalLevel = allowed ? newLevel + 1 : newLevel

    // Update bucket state
    const updateMulti = redis.multi()
    updateMulti.hset(key, 'lastUpdate', now.toString())
    updateMulti.hset(key, 'level', finalLevel.toString())
    updateMulti.pexpire(key, 60000) // Expire after 1 minute of inactivity
    await updateMulti.exec()

    // Calculate estimated wait time if not allowed
    const estimatedWaitMs = allowed ? 0 : Math.ceil(1000 / leakRate)

    return {
      allowed,
      currentLevel: finalLevel,
      capacity: options.bucketCapacity,
      estimatedWaitMs,
    }
  }

  /**
   * Reset a specific bucket
   */
  async function reset(bucketName: string): Promise<void> {
    const key = `${prefix}${bucketName}`
    await redis.del(key)
  }

  return {
    isAllowed,
    reset,
  }
}
