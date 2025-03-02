# Base image with Bun - use platform-specific image
FROM --platform=${BUILDPLATFORM:-linux/amd64} oven/bun:canary-alpine AS alpine
RUN apk update
RUN apk add --no-cache libc6-compat tree nodejs npm

# Setup pnpm and turbo
RUN npm install pnpm turbo --global
RUN pnpm config set store-dir ~/.pnpm-store

# Prune projects
FROM alpine AS pruner
ARG PROJECT=tanstack

WORKDIR /app
COPY . .
RUN turbo prune --scope=${PROJECT} --docker

# Build the project
FROM alpine AS builder
ARG PROJECT=tanstack

WORKDIR /app

# Copy lockfile and package.json's of isolated subworkspace
COPY --from=pruner /app/out/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=pruner /app/out/pnpm-workspace.yaml ./pnpm-workspace.yaml
COPY --from=pruner /app/out/json/ .

# First install the dependencies (as they change less often)
RUN --mount=type=cache,id=pnpm,target=~/.pnpm-store pnpm install

# Copy source code of isolated subworkspace
COPY --from=pruner /app/out/full/ .

RUN echo "Contents before build:" && ls -la apps/tanstack/
RUN echo "\nturbo.json contents:" && cat turbo.json

# Build the app
WORKDIR /app

RUN pnpm --filter ${PROJECT} build
RUN echo "Contents after build:" && ls -la apps/tanstack/
RUN tree apps/tanstack/

# Return to app root
WORKDIR /app

RUN --mount=type=cache,id=pnpm,target=~/.pnpm-store pnpm prune --prod --no-optional

# Final image - use platform-specific image
FROM --platform=${TARGETPLATFORM:-linux/amd64} oven/bun:canary-alpine AS runner
ARG PROJECT=tanstack

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nodejs
USER nodejs

WORKDIR /app

# Copy only the built files and dependencies
COPY --from=builder --chown=nodejs:nodejs /app/apps/tanstack/.output ./.output
COPY --from=builder --chown=nodejs:nodejs /app/apps/tanstack/package.json ./package.json
# COPY --from=builder /app/.prod_modules ./node_modules


ARG PORT=3004
ENV PORT=${PORT}
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
EXPOSE ${PORT}

# Start the server
CMD ["bun", "run", "start"]
