import { join } from 'path'
import { serve } from 'bun'

const distPath = './dist'

serve({
  port: 3001,
  fetch(req) {
    const url = new URL(req.url)
    const path = url.pathname === '/' ? '/index.html' : url.pathname

    try {
      const file = Bun.file(join(distPath, path))
      return new Response(file)
    } catch (e) {
      return new Response('Not Found', { status: 404 })
    }
  },
})

console.log('Server running at http://localhost:3001')
