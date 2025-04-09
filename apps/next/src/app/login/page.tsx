'use client'

import { useEffect, useState } from 'react'

const API_BASE = 'http://localhost:3000'

export default function AuthTest() {
  const [csrfToken, setCsrfToken] = useState('')
  const [status, setStatus] = useState('')

  useEffect(() => {
    // Get CSRF token on page load
    fetch(`${API_BASE}/api/auth/csrf-token`)
      .then((res) => res.json())
      .then((data) => setCsrfToken(data.csrf_token))
  }, [])

  const handleLogin = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/auth/login/google`, {
        method: 'POST',
        headers: {
          'X-CSRF-Token': csrfToken,
          'Content-Type': 'application/json',
        },
      })
      const data = await res.json()
      setStatus(`Login response: ${JSON.stringify(data)}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  const handleLogout = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/auth/logout`, {
        method: 'POST',
        headers: {
          'X-CSRF-Token': csrfToken,
          'Content-Type': 'application/json',
        },
      })
      const data = await res.json()
      setStatus(`Logout response: ${JSON.stringify(data)}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  const handleSessionsPost = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/auth/sessionsPost`, {
        method: 'POST',
        headers: {
          'X-CSRF-Token': csrfToken,
          'Content-Type': 'application/json',
        },
      })
      const data = await res.json()
      setStatus(`Sessions response: ${JSON.stringify(data)}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  const handleSessionsGet = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/auth/sessions`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await res.json()
      setStatus(`Sessions GET response: ${JSON.stringify(data)}`)
    } catch (error: unknown) {
      console.error(error)
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }
  const handleGetStatus = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/auth/status`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await res.json()
      setStatus(`Status response: ${JSON.stringify(data)}`)
    } catch (error: unknown) {
      setStatus(
        `Error: ${error instanceof Error ? error.message : 'Unknown error'}`
      )
    }
  }

  const handleClearToken = () => {
    setCsrfToken('')
    setStatus('CSRF token cleared')
  }

  return (
    <div className="container p-4">
      <h1 className="mb-4 text-2xl">Auth Test Page</h1>

      <div className="space-y-4">
        <div>
          <p className="text-sm text-gray-600">CSRF Token:</p>
          <code className="block rounded bg-gray-100 p-2">{csrfToken}</code>
        </div>

        <div className="space-x-4">
          <button
            onClick={handleLogin}
            className="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
          >
            Test Login
          </button>

          <button
            onClick={handleLogout}
            className="rounded bg-red-500 px-4 py-2 text-white hover:bg-red-600"
          >
            Test Logout
          </button>

          <button
            onClick={handleSessionsPost}
            className="rounded bg-green-500 px-4 py-2 text-white hover:bg-green-600"
          >
            Test Sessions POST
          </button>

          <button
            onClick={handleSessionsGet}
            className="rounded bg-purple-500 px-4 py-2 text-white hover:bg-purple-600"
          >
            Test Sessions GET
          </button>

          <button
            onClick={handleGetStatus}
            className="rounded bg-yellow-500 px-4 py-2 text-white hover:bg-yellow-600"
          >
            Test Status
          </button>

          <button
            onClick={handleClearToken}
            className="rounded bg-gray-500 px-4 py-2 text-white hover:bg-gray-600"
          >
            Clear Token
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
