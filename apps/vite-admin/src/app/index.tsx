import './styles.css'

import { SelectDemo } from '@repo/ui/components/styled'
import { CounterButton } from '@repo/ui/counter-button'
import { Link } from '@repo/ui/link'

function App() {
  return (
    <div className="container">
      <h1 className="title">
        Admin <br />
        <span>Kitchen Sink</span>
      </h1>
      <CounterButton />
      <SelectDemo />
      <p className="description">
        Built With{' '}
        <Link href="https://turbo.build/repo" newTab>
          Turborepo
        </Link>
        {' & '}
        <Link href="https://vitejs.dev/" newTab>
          Vite
        </Link>
      </p>
    </div>
  )
}

export default App
