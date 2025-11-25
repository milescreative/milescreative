import { SQL } from 'bun'

console.log('Initializing database connection with Bun SQL...')

const config = {
  hostname: '192.168.1.50',
  port: 5433,
  database: 'postgres',
  username: 'postgres',
  password: process.env.PG_PASSWORD,
}

const sql = new SQL(config)

sql.options.onconnect = (): void => {
  console.log('Connected to database')
}

sql.options.onclose = (): void => {
  console.log('Connection closed')
}

export default sql
