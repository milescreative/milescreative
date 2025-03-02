import * as fs from 'fs'
import { glob } from 'glob'

const WORKSPACE_SETTINGS = {
  repo_name: '@repo', // CHANGE BY RUNNING PACKAGE.JSON SCRIPT
  preferred_runtime: 'bun',
}

interface PackageJson {
  name?: string
  dependencies?: Record<string, string>
  devDependencies?: Record<string, string>
  [key: string]: any
}

interface TsConfig {
  compilerOptions?: {
    paths?: Record<string, string[]>
  }
  extends?: string
  [key: string]: any
}

function updateWorkspaceSettingsFile(newPrefix: string) {
  const filePath = __dirname + '/workspace.ts'
  const content = fs.readFileSync(filePath, 'utf-8')
  const newContent = content.replace(
    /repo_name:\s*['"]@[^'"]+['"]/,
    `repo_name: '${newPrefix}'`
  )
  fs.writeFileSync(filePath, newContent)
}

export function renameWorkspacePackages(newPrefix: string): void {
  if (!newPrefix.startsWith('@')) {
    throw new Error('New prefix must start with @')
  }

  // Remove trailing slash if present
  newPrefix = newPrefix.endsWith('/') ? newPrefix.slice(0, -1) : newPrefix
  const oldPrefix = WORKSPACE_SETTINGS.repo_name

  // Step 1: Remove all node_modules
  console.log('Removing node_modules directories...')
  try {
    const nodeModulesPaths = glob.sync('**/node_modules', {
      ignore: ['**/node_modules/**/node_modules/**'],
    })

    for (const modulePath of nodeModulesPaths) {
      console.log(`Removing ${modulePath}...`)
      fs.rmSync(modulePath, { recursive: true, force: true })
    }
  } catch (error) {
    console.error('Error removing node_modules:', error)
    throw error
  }

  // Step 2: Find all package.json files
  const packageJsonFiles = glob.sync('**/package.json', {
    ignore: ['**/node_modules/**', '**/dist/**', '**/build/**'],
  })

  // Step 3: Process each package.json
  for (const filePath of packageJsonFiles) {
    try {
      const content = fs.readFileSync(filePath, 'utf-8')
      const packageJson: PackageJson = JSON.parse(content)
      let modified = false

      // Update package name if it starts with the old prefix
      if (packageJson.name?.startsWith(oldPrefix)) {
        packageJson.name = packageJson.name.replace(oldPrefix, newPrefix)
        modified = true
      }

      // Update dependencies
      if (packageJson.dependencies) {
        for (const [dep, version] of Object.entries(packageJson.dependencies)) {
          if (dep.startsWith(oldPrefix)) {
            const newDep = dep.replace(oldPrefix, newPrefix)
            delete packageJson.dependencies[dep]
            packageJson.dependencies[newDep] = version
            modified = true
          }
        }
      }

      // Update devDependencies
      if (packageJson.devDependencies) {
        for (const [dep, version] of Object.entries(
          packageJson.devDependencies
        )) {
          if (dep.startsWith(oldPrefix)) {
            const newDep = dep.replace(oldPrefix, newPrefix)
            delete packageJson.devDependencies[dep]
            packageJson.devDependencies[newDep] = version
            modified = true
          }
        }
      }

      // Save changes if modified
      if (modified) {
        fs.writeFileSync(filePath, JSON.stringify(packageJson, null, 2) + '\n')
        console.log(`Updated ${filePath}`)
      }
    } catch (error) {
      console.error(`Error processing ${filePath}:`, error)
      throw error
    }
  }

  // Step 4: Update tsconfig paths
  const tsconfigFiles = glob.sync('**/tsconfig.json', {
    ignore: ['**/node_modules/**', '**/dist/**', '**/build/**'],
  })

  for (const filePath of tsconfigFiles) {
    try {
      const content = fs.readFileSync(filePath, 'utf-8')
      const tsconfig: TsConfig = JSON.parse(content)
      let modified = false

      // Update extends field if it exists and contains the old prefix
      if (
        typeof tsconfig.extends === 'string' &&
        tsconfig.extends.startsWith(oldPrefix)
      ) {
        tsconfig.extends = tsconfig.extends.replace(oldPrefix, newPrefix)
        modified = true
      }

      // Update paths
      if (tsconfig.compilerOptions?.paths) {
        const newPaths: Record<string, string[]> = {}
        for (const [key, value] of Object.entries(
          tsconfig.compilerOptions.paths
        )) {
          if (key.startsWith(oldPrefix)) {
            const newKey = key.replace(oldPrefix, newPrefix)
            newPaths[newKey] = value
            modified = true
          } else {
            newPaths[key] = value
          }
        }
        tsconfig.compilerOptions.paths = newPaths
      }

      if (modified) {
        fs.writeFileSync(filePath, JSON.stringify(tsconfig, null, 2) + '\n')
        console.log(`Updated ${filePath}`)
      }
    } catch (error) {
      console.error(`Error processing ${filePath}:`, error)
      throw error
    }
  }

  // Step 5: Update imports in source files and configs
  const sourceFiles = glob.sync('**/*.{ts,tsx,css,json,js,mjs,cjs,md,astro}', {
    ignore: [
      '**/node_modules/**',
      '**/dist/**',
      '**/build/**',
      '**/package.json',
      '**/tsconfig.json',
      '**/pnpm-lock.yaml',
    ],
  })

  for (const filePath of sourceFiles) {
    try {
      const content = fs.readFileSync(filePath, 'utf-8')
      let newContent = content

      // Handle different file types appropriately
      if (filePath.endsWith('.css')) {
        // CSS imports
        newContent = content.replace(
          new RegExp(`@import\\s+['"]${oldPrefix}/([^'"]+)['"]`, 'g'),
          `@import '${newPrefix}/$1'`
        )
      } else if (filePath.endsWith('.md')) {
        // Markdown files - handle both inline code and links
        newContent = content.replace(
          new RegExp(`${oldPrefix}/([^\\s'"\`\\)]+)`, 'g'),
          `${newPrefix}/$1`
        )
      } else {
        // JS/TS/Astro files - handle imports and requires
        const importRegex = new RegExp(`['"]${oldPrefix}/([^'"]+)['"]`, 'g')
        newContent = content.replace(importRegex, `'${newPrefix}/$1'`)
      }

      if (content !== newContent) {
        fs.writeFileSync(filePath, newContent)
        console.log(`Updated references in ${filePath}`)
      }
    } catch (error) {
      console.error(`Error processing ${filePath}:`, error)
      throw error
    }
  }

  // Update workspace settings
  WORKSPACE_SETTINGS.repo_name = newPrefix
  updateWorkspaceSettingsFile(newPrefix)

  console.log('Package renaming completed successfully.')
  console.log(
    "Please run your package manager's install command to rebuild node_modules."
  )
}
