import { createInterface } from 'readline/promises'

import { renameWorkspacePackages } from './workspace'

async function main() {
  const rl = createInterface({
    input: process.stdin,
    output: process.stdout,
  })

  try {
    const newPrefix = await rl.question(
      'Enter new workspace prefix (e.g. @myorg): '
    )

    if (!newPrefix || !newPrefix.startsWith('@')) {
      console.error('Error: Prefix must start with @ and not be empty')
      process.exit(1)
    }

    renameWorkspacePackages(newPrefix)
  } catch (error) {
    console.error('Error:', error)
    process.exit(1)
  } finally {
    rl.close()
  }
}

main()
