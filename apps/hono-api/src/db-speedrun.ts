import { performance } from 'node:perf_hooks'
import { Database } from 'bun:sqlite'
import { Redis } from 'ioredis'

const redis = new Redis(4545, 'milescreative-s1', {
  username: 'default',
  password: 'TuryEACRshgH2TzF10uXqFZA1qZTHDy2ESO02AZudzJllzpvmKdrEbrxgSizlWAT',
})

const db = new Database('rate-limit.db', { create: true })

// Setup SQLite table
db.run(`
  CREATE TABLE IF NOT EXISTS rate_limits (
    key TEXT PRIMARY KEY,
    requests INTEGER,
    timestamp INTEGER
  )
`)

async function simulateConcurrentRedis(
  totalRequests: number,
  concurrentUsers: number
) {
  const start = performance.now()
  const requestsPerUser = Math.floor(totalRequests / concurrentUsers)

  const userPromises = Array.from(
    { length: concurrentUsers },
    async (_, userIndex) => {
      const key = `test-user-${userIndex}`
      const requests = Array.from({ length: requestsPerUser }, () =>
        redis.multi().incr(`${key}:count`).pexpire(`${key}:count`, 60000).exec()
      )
      await Promise.all(requests)
    }
  )

  await Promise.all(userPromises)
  const duration = (performance.now() - start) / 1000 // Convert to seconds

  return {
    totalTime: duration,
    requestsPerSecond: totalRequests / duration,
    avgLatencyMs: (duration * 1000) / totalRequests,
  }
}

function simulateConcurrentSQLite(
  totalRequests: number,
  concurrentUsers: number
) {
  const start = performance.now()
  const requestsPerUser = Math.floor(totalRequests / concurrentUsers)

  // SQLite statement preparation (reusable)
  const stmt = db.prepare(`
    INSERT INTO rate_limits (key, requests, timestamp)
    VALUES (?, 1, ?)
    ON CONFLICT(key) DO UPDATE SET
    requests = requests + 1,
    timestamp = ?
  `)

  // Since SQLite is single-threaded, we'll simulate concurrent access
  for (let userIndex = 0; userIndex < concurrentUsers; userIndex++) {
    const key = `test-user-${userIndex}`
    for (let i = 0; i < requestsPerUser; i++) {
      const now = Date.now()
      stmt.run(key, now, now)
    }
  }

  const duration = (performance.now() - start) / 1000 // Convert to seconds

  return {
    totalTime: duration,
    requestsPerSecond: totalRequests / duration,
    avgLatencyMs: (duration * 1000) / totalRequests,
  }
}

async function runBenchmark() {
  const scenarios = [{ requests: 10000, users: 100 }]

  console.log('Starting benchmark...\n')

  for (const { requests, users } of scenarios) {
    console.log(
      `\nScenario: ${requests} requests with ${users} concurrent users`
    )
    console.log('----------------------------------------')

    // Redis benchmark
    const redisResults = await simulateConcurrentRedis(requests, users)
    console.log('\nRedis Results:')
    console.log(`Total time: ${redisResults.totalTime.toFixed(2)}s`)
    console.log(`Requests/second: ${redisResults.requestsPerSecond.toFixed(2)}`)
    console.log(`Avg latency: ${redisResults.avgLatencyMs.toFixed(2)}ms`)

    // SQLite benchmark
    const sqliteResults = simulateConcurrentSQLite(requests, users)
    console.log('\nSQLite Results:')
    console.log(`Total time: ${sqliteResults.totalTime.toFixed(2)}s`)
    console.log(
      `Requests/second: ${sqliteResults.requestsPerSecond.toFixed(2)}`
    )
    console.log(`Avg latency: ${sqliteResults.avgLatencyMs.toFixed(2)}ms`)
  }

  // Cleanup
  await redis.quit()
  db.close()
}

runBenchmark()
