import type { Redis } from 'ioredis'

interface RateLimiterConfig {
  redis: Redis
  appPrefix: string
}

interface UserRateLimiterOptions {
  /**
   * Maximum number of requests allowed in the time window
   */
  maxRequests: number

  /**
   * Time window in seconds
   */
  windowSeconds: number
}

interface RateLimitResult {
  /**
   * Whether the request is allowed based on rate limits
   */
  allowed: boolean

  /**
   * Number of requests remaining in the current time window
   */
  remaining: number

  /**
   * When the rate limit will reset
   */
  resetAt: Date

  /**
   * Total number of requests allowed in the time window
   */
  total: number
}

/**
 * User rate limiter functions
 */
export interface UserRateLimiter {
  /**
   * Check if a user has exceeded their rate limit
   * @param userId Unique identifier for the user
   * @returns Object containing limit information and whether the request is allowed
   */
  isAllowed(userId: string): Promise<RateLimitResult>

  /**
   * Reset rate limit for a specific user
   * @param userId Unique identifier for the user
   */
  reset(userId: string): Promise<void>
}

/**
 * Create a rate limiter function that checks if a user has exceeded their rate limit
 *
 * @param config Configuration with Redis instance and app prefix
 * @param options Rate limiting options with defaults
 * @returns A function that checks if a user is allowed to make requests
 */
export function createUserRateLimiter(
  config: RateLimiterConfig,
  options: UserRateLimiterOptions = { maxRequests: 10, windowSeconds: 10 }
): UserRateLimiter {
  const redis = config.redis
  const prefix = `${config.appPrefix}:user-ratelimit:`

  /**
   * Check if a user has exceeded their rate limit
   */
  async function isAllowed(userId: string): Promise<RateLimitResult> {
    const key = `${prefix}${userId}`
    const now = Date.now()
    const windowMs = options.windowSeconds * 1000
    const windowStart = now - windowMs

    // Remove expired timestamps and add current request timestamp
    const multi = redis.multi()
    multi.zremrangebyscore(key, 0, windowStart)
    multi.zadd(key, now, `${now}`)
    multi.zrange(key, 0, -1)
    multi.pexpire(key, windowMs)

    const results = await multi.exec()
    if (!results) {
      throw new Error('Redis operation failed')
    }

    // Type assertion for Redis response
    // Make sure we properly handle the Redis response, which could be undefined
    const requestTimestamps = (results[2]?.[1] ?? []) as string[]
    const count = requestTimestamps.length
    const resetAt = new Date(
      parseInt(requestTimestamps[0] || `${now}`, 10) + windowMs
    )

    return {
      allowed: count <= options.maxRequests,
      remaining: Math.max(0, options.maxRequests - count),
      resetAt,
      total: options.maxRequests,
    }
  }

  /**
   * Reset rate limit for a specific user
   */
  async function reset(userId: string): Promise<void> {
    const key = `${prefix}${userId}`
    await redis.del(key)
  }

  return {
    isAllowed,
    reset,
  }
}
