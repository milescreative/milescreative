// import { useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import {  useState } from 'react'
import csrfStore from '../data/store'
import SomeTest from '../components/ui/some-test'
import * as cache from '../data/cache'
import * as myAtom from '../data/atom'





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
        window.location.href = 'http://localhost:3000/api/auth/login?redirect_url=http://localhost:3001/auth'
        return
      }

      const response = await fetch(`http://localhost:3000/api/auth/${endpoint}`, {
        method,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await response.json()
      setResponse(JSON.stringify(data, null, 2))
    } catch (error) {
      setResponse(`Error: ${error instanceof Error ? error.message : String(error)}`)
    } finally {
      setLoading(false)
    }
  }



  const handleCSRFProtectedApiCall = async (endpoint: string, method: string = 'POST') => {
    // const csrfToken = csrfStore.get()
    // const csrfToken2 = cache.get()
    const csrfToken3 = myAtom.csrfStore.get()

    myAtom.csrfStore.off();
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

      const response = await fetch(`http://localhost:3000/api/auth/${endpoint}`, {
        method,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken3 || '',
        },
      })
      const data = await response.json()
      setResponse(JSON.stringify(data, null, 2))
      } catch (error) {
      setResponse(`Error: ${error instanceof Error ? error.message : String(error)}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Auth API Test Page</h1>

      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="space-y-4">
          <button
            onClick={() => handleApiCall('login')}
            className="w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
            disabled={loading}
          >
            Login
          </button>

          <button
            onClick={() => handleApiCall('validate')}
            className="w-full px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:opacity-50"
            disabled={loading}
          >
            Validate Session
          </button>

          <button
            onClick={() => handleApiCall('refresh')}
            className="w-full px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600 disabled:opacity-50"
            disabled={loading}
          >
            Refresh Token
          </button>

          <button
            onClick={() => handleApiCall('logout')}
            className="w-full px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            Logout
          </button>

          <button
            onClick={() => handleApiCall('csrf')}
            className="w-full px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            CSRF
          </button>

          <button
            onClick={() => handleCSRFProtectedApiCall('csrf-protected')}
            className="w-full px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:opacity-50"
            disabled={loading}
          >
            CSRF Protected
          </button>
        </div>

        <div className="bg-gray-100 p-4 rounded">
          <h2 className="text-lg font-semibold mb-2">Response:</h2>
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
        <h2 className="text-lg font-semibold mb-2">User Management</h2>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-4">
            <button
              onClick={() => handleApiCall('user?user_id=test')}
              className="w-full px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 disabled:opacity-50"
              disabled={loading}
            >
              Get User
            </button>

            <button
              onClick={() => handleApiCall('user/sessions?user_id=test')}
              className="w-full px-4 py-2 bg-indigo-500 text-white rounded hover:bg-indigo-600 disabled:opacity-50"
              disabled={loading}
            >
              Get User Sessions
            </button>

            <button
              onClick={() => handleApiCall('user?user_id=test&name=Test&email=test@test.com', 'PUT')}
              className="w-full px-4 py-2 bg-pink-500 text-white rounded hover:bg-pink-600 disabled:opacity-50"
              disabled={loading}
            >
              Update User
            </button>

            <button
              onClick={() => handleApiCall('user?user_id=test', 'DELETE')}
              className="w-full px-4 py-2 bg-orange-500 text-white rounded hover:bg-orange-600 disabled:opacity-50"
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
