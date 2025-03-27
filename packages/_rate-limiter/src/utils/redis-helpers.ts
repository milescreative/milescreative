import { ChainableCommander } from 'ioredis'
import { z } from 'zod'

export async function rtxn_exec<Schemas extends z.ZodType[]>(
  txn: ChainableCommander,
  schemas: [...Schemas]
): Promise<{ [K in keyof Schemas]: z.infer<Schemas[K]> }> {
  const results = await txn.exec()

  if (results === null) {
    throw new Error('Failed to get token bucket state')
  }

  if (results.length !== schemas.length) {
    throw new Error('Number of results does not match number of schemas')
  }

  return results.map(([error, value], index) => {
    if (error) {
      throw error
    }
    return schemas[index].parse(value)
  }) as { [K in keyof Schemas]: z.infer<Schemas[K]> }
}
