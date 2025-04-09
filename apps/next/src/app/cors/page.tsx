'use client'

import { useState } from 'react'

const GO_API_BASE = 'http://localhost:8080'

export default function CorsTest() {
  const [status, setStatus] = useState('')

  const handleTestCors = async () => {
    try {
      const res = await fetch(`${GO_API_BASE}/test-cors`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await res.text()
      setStatus(`CORS Test Response: ${data}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  const handleTestPost = async () => {
    try {
      const res = await fetch(`${GO_API_BASE}/postTest`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await res.text()
      setStatus(`POST Test Response: ${data}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  return (
    <div className="container p-4">
      <h1 className="mb-4 text-2xl">CORS Test Page</h1>

      <div className="space-y-4">
        <div className="space-x-4">
          <button
            onClick={handleTestCors}
            className="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
          >
            Test CORS GET
          </button>

          <button
            onClick={handleTestPost}
            className="rounded bg-green-500 px-4 py-2 text-white hover:bg-green-600"
          >
            Test CORS POST
          </button>
        </div>

        {status && (
          <div className="mt-4">
            <p className="text-sm text-gray-600">Response:</p>
            <pre className="mt-1 rounded bg-gray-100 p-2">{status}</pre>
          </div>
        )}
      </div>
    </div>
  )
}
