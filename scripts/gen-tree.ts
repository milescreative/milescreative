import { execSync } from 'child_process'
import { appendFileSync, readFileSync } from 'fs'
import { glob } from 'glob'

// Initialize output file with tree structure
console.log('Generating project information...')
execSync('tree -L 5 -I "node_modules" > ../project-info.txt')

// Function to append package.json content to file
function appendPackageJson(path: string) {
  try {
    const content = readFileSync(path, 'utf-8')
    const separator = `\n\n=== Package.json for ${path} ===\n\n`
    appendFileSync('project-info.txt', separator + content)
  } catch (error) {
    console.error(`Error reading ${path}:`, error)
  }
}

// Append root package.json
appendPackageJson('../package.json')

// Append all package.json files in apps directory
const appPackageJsons = glob.sync('../apps/*/package.json')
appPackageJsons.forEach(appendPackageJson)

console.log('Project information has been written to project-info.txt')
