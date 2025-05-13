import { createFileRoute } from '@tanstack/react-router'
import csrfStore from '../data/store'
import { useState } from 'react'
export const Route = createFileRoute('/csrf-again')({
  component: RouteComponent,
})

function RouteComponent() {
  const csrfToken = csrfStore.get()
  const [response, setResponse] = useState<string>('')

  const handleClick = async () => {
    const response = await fetch('http://localhost:3000/api/auth/csrf-protected', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-TOKEN': csrfToken,

        },
        body: JSON.stringify({ some: 'data' }),
      });
    setResponse(await response.json())
  }
  return <div>Hello "/csrf-again"!
    <button onClick={handleClick}>Click me</button>
    <p>{`response: ${response}`}</p>
  </div>
}
