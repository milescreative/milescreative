import sql from './db'

export const test_result = async () => {
  const result = await sql`SELECT * FROM test`
  console.log(result)
  return result
}
