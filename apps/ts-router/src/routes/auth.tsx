// import { useQuery } from '@tanstack/react-query';
import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { safeFetch } from '../../../../packages/backend-sdk/src'
import SomeTest from '../components/ui/some-test'
import { csrfStore } from '../data/atom'

export const Route = createFileRoute('/auth')({
  component: RouteComponent,
})

function RouteComponent() {
  // const csrfQuery = useQuery({
  //   queryKey: ['csrf_token'],
  //   queryFn: async () => {
  //     const resp = await fetch('http://localhost:3000/api/auth/csrf', {
  //       credentials: 'include',
  //       method: 'GET',
  //     })
  //     for (const [key, value] of resp.headers.entries()) {
  //       console.log('header', key, value)
  //     }
  //     const token = resp.headers.get('X-Csrf-Token')
  //     if (!token) {
  //       throw new Error('No CSRF token found')
  //     }
  //     return token
  //   },
  //   refetchOnWindowFocus: false,
  //   refetchOnMount: false,
  //   refetchOnReconnect: false,
  //   refetchInterval: false,
  //   staleTime: Infinity,
  //   gcTime: Infinity,
  //   retry: true

  // })

  const [response, setResponse] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  const handleApiCall = async (endpoint: string, method: string = 'GET') => {
    setLoading(true)
    try {
      // Special handling for login since it's a redirect flow
      if (endpoint === 'login') {
        window.location.href =
          'http://localhost:3000/api/auth/login/google?redirect_url=http://localhost:3001/auth'
        return
      }

      const response = await fetch(
        `http://localhost:3000/api/auth/${endpoint}`,
        {
          method,
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      )
      const data = await response.json()
      setResponse(JSON.stringify(data, null, 2))
    } catch (error) {
      setResponse(
        `Error: ${error instanceof Error ? error.message : String(error)}`,
      )
    } finally {
      setLoading(false)
    }
  }

  const handleCSRFProtectedApiCall = async (endpoint: string) => {
    // const csrfToken = csrfStore.get()
    // const csrfToken2 = cache.get()
    const csrfToken3 = csrfStore.get()

    csrfStore.off()
    setLoading(true)
    try {
      // if (!csrfQuery.data) {
      //   await csrfQuery.refetch()
      //   if (!csrfQuery.data) {
      //     setResponse('No CSRF token found')
      //     return
      //   }
      // }

      // console.log('csrfToken', csrfToken)

      const response = await safeFetch(
        `http://localhost:3000/api/auth/${endpoint}`,
        {},
        () => csrfToken3,
      )
      const data = await response.json()
      setResponse(JSON.stringify(data, null, 2))
    } catch (error) {
      setResponse(
        `Error: ${error instanceof Error ? error.message : String(error)}`,
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto max-w-4xl p-6">
      <h1 className="mb-6 text-2xl font-bold">Auth API Test Page</h1>

      <div className="mb-6 grid grid-cols-2 gap-4">
        <div className="space-y-4">
          <button
            onClick={() => handleApiCall('login')}
            className="w-full rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600 disabled:opacity-50"
            disabled={loading}
          >
            Login
          </button>

          <button
            onClick={() => handleApiCall('validate')}
            className="w-full rounded bg-green-500 px-4 py-2 text-white hover:bg-green-600 disabled:opacity-50"
            disabled={loading}
          >
            Validate Session
          </button>

          <button
            onClick={() => handleApiCall('refresh', 'POST')}
            className="w-full rounded bg-yellow-500 px-4 py-2 text-white hover:bg-yellow-600 disabled:opacity-50"
            disabled={loading}
          >
            Refresh Token
          </button>

          <button
            onClick={() => handleApiCall('logout', 'POST')}
            className="w-full rounded bg-red-500 px-4 py-2 text-white hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            Logout
          </button>

          <button
            onClick={() => handleApiCall('csrf')}
            className="w-full rounded bg-red-500 px-4 py-2 text-white hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            CSRF
          </button>

          <button
            onClick={() => handleCSRFProtectedApiCall('csrf-protected')}
            className="w-full rounded bg-red-500 px-4 py-2 text-white hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            CSRF Protected
          </button>
        </div>

        <div className="rounded bg-gray-100 p-4">
          <h2 className="mb-2 text-lg font-semibold">Response:</h2>
          {loading ? (
            <div className="text-gray-500">Loading...</div>
          ) : (
            <pre className="whitespace-pre-wrap break-words">
              {response || 'No response yet'}
            </pre>
          )}
        </div>
      </div>

      <div className="mt-6">
        <h2 className="mb-2 text-lg font-semibold">User Management</h2>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-4">
            <button
              onClick={() => handleApiCall('user?user_id=test')}
              className="w-full rounded bg-purple-500 px-4 py-2 text-white hover:bg-purple-600 disabled:opacity-50"
              disabled={loading}
            >
              Get User
            </button>

            <button
              onClick={() => handleApiCall('user/sessions?user_id=test')}
              className="w-full rounded bg-indigo-500 px-4 py-2 text-white hover:bg-indigo-600 disabled:opacity-50"
              disabled={loading}
            >
              Get User Sessions
            </button>

            <button
              onClick={() =>
                handleApiCall(
                  'user?user_id=test&name=Test&email=test@test.com',
                  'PUT',
                )
              }
              className="w-full rounded bg-pink-500 px-4 py-2 text-white hover:bg-pink-600 disabled:opacity-50"
              disabled={loading}
            >
              Update User
            </button>

            <button
              onClick={() => handleApiCall('user?user_id=test', 'DELETE')}
              className="w-full rounded bg-orange-500 px-4 py-2 text-white hover:bg-orange-600 disabled:opacity-50"
              disabled={loading}
            >
              Delete User
            </button>
          </div>
        </div>
      </div>
      <SomeTest />
    </div>
  )
}
