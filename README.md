# Modern Monorepo Boilerplate

A modern, high-performance monorepo setup using Turborepo, pnpm, and Bun.

## 🚀 Features

- **Package Manager**: [pnpm](https://pnpm.io/) for fast, disk-space efficient package management
- **Build System**: [Turborepo](https://turbo.build/) for optimized build pipelines
- **Runtime**: [Bun](https://bun.sh/) for blazing fast JavaScript runtime and tooling
  - All apps and packages are configured to use Bun as the default runtime
  - Leverages Bun's built-in bundler, test runner, and package manager capabilities
- **Type Safety**: Full TypeScript support across all packages
- **Code Quality**: ESLint and Prettier configuration with import sorting
- **Docker Support**: Ready-to-use Docker configuration
  - Production-ready Dockerfiles for each application
  - Optimized for [Coolify](https://coolify.io/) deployment
  - Multi-stage builds with best practices

## 📦 Prerequisites

- Node.js (LTS version recommended)
- pnpm 10.5.2 or higher
- Bun 1.2.2 or higher

## 🛠️ Getting Started

1. **Clone the repository**

```bash
git clone <your-repo-url>
cd <repo-name>
```

2. **Install dependencies**

```bash
pnpm install
```

3. **Start development**

```bash
pnpm dev
```

## 📚 Available Scripts

- `pnpm build` - Build all packages
- `pnpm dev` - Start development mode
- `pnpm lint` - Lint all packages
- `pnpm test` - Run tests
- `pnpm clean` - Clean build artifacts
- `pnpm format` - Format code with Prettier
- `pnpm check` - Run sherif checks
- `pnpm docker:build` - Build Docker containers
- `pnpm docker:up` - Start Docker environment
- `pnpm @up` - Update dependencies recursively

## 🏗️ Project Structure

```
.
├── apps/                    # Application packages
│   ├── vite-admin/         # Vite-powered admin dashboard
│   ├── astro/             # Astro.js application
│   ├── hono-api/          # Hono-based API
│   ├── next/              # Next.js application
│   └── tanstack/          # TanStack application
├── packages/               # Shared packages
│   ├── ui/                # Shared UI components
│   ├── config-eslint/     # ESLint configurations
│   ├── db/                # Database utilities
│   ├── design/            # Design system
│   ├── logger/            # Logging utilities
│   └── config-typescript/ # TypeScript configurations
├── dockerfiles/           # Docker configuration files
├── docker-compose.yml     # Docker compose configuration
├── turbo.json            # Turborepo configuration
└── package.json          # Root package.json
```

## 🔧 Development Tools

- **TypeScript**: Version 5.7.3
- **ESLint**: Version 9.20.0
- **Prettier**: Version 3.5.2
  - With plugins for import sorting, package.json formatting, and Tailwind CSS

## 🐳 Docker Support

The project includes Docker support for containerized development and deployment, with pre-configured Dockerfiles optimized for [Coolify](https://coolify.io/) deployment:

```bash
# Build Docker containers
pnpm docker:build

# Start Docker environment
pnpm docker:up
```

### 🚀 Coolify Deployment

Each application in the `apps` directory includes a production-ready Dockerfile configured for seamless deployment to Coolify:

- Optimized multi-stage builds
- Built-in health checks
- Environment variable handling
- Automatic port configuration
- Volume persistence setup
- Resource optimization

To deploy to Coolify:

1. Connect your Git repository to Coolify
2. Select the application directory (e.g., `apps/next` or `apps/hono-api`)
3. The pre-configured Dockerfile will be automatically detected
4. Configure your environment variables
5. Deploy!

## 📝 Code Style

This project uses Prettier with the following plugins:

- `@ianvs/prettier-plugin-sort-imports`
- `prettier-plugin-packagejson`
- `prettier-plugin-tailwindcss`

## 🤝 Contributing

1. Create a new branch
2. Make your changes
3. Submit a pull request

## 📄 License

[MIT](LICENSE)

WIP
