# Base image with Bun - use platform-specific image
FROM --platform=${BUILDPLATFORM:-linux/amd64} oven/bun:canary-alpine AS alpine
RUN apk update
RUN apk add --no-cache libc6-compat tree nodejs npm

# Setup pnpm and turbo
RUN npm install pnpm turbo --global
RUN pnpm config set store-dir ~/.pnpm-store

# Prune projects
FROM alpine AS pruner
ARG PROJECT=next

WORKDIR /app
COPY . .
RUN turbo prune --scope=${PROJECT} --docker

# Build the project
FROM alpine AS builder
ARG PROJECT=next

WORKDIR /app

# Copy lockfile and package.json's of isolated subworkspace
COPY --from=pruner /app/out/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=pruner /app/out/pnpm-workspace.yaml ./pnpm-workspace.yaml
COPY --from=pruner /app/out/json/ .

# First install the dependencies (as they change less often)
RUN --mount=type=cache,id=pnpm,target=~/.pnpm-store pnpm install --frozen-lockfile

# Copy source code of isolated subworkspace
COPY --from=pruner /app/out/full/ .

RUN turbo build --filter=${PROJECT}
RUN --mount=type=cache,id=pnpm,target=~/.pnpm-store pnpm prune --prod --no-optional
RUN rm -rf ./**/*/src
# Add debug command to see what files we have after build


# Final image - use platform-specific image
FROM --platform=${TARGETPLATFORM:-linux/amd64} oven/bun:canary-alpine AS runner
ARG PROJECT=next

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nodejs
USER nodejs

WORKDIR /app

# Copy standalone directory and its contents
COPY --from=builder --chown=nodejs:nodejs /app/apps/next/.next/standalone/. ./
COPY --from=builder --chown=nodejs:nodejs /app/apps/next/.next/static ./apps/next/.next/static
COPY --from=builder --chown=nodejs:nodejs /app/apps/next/public ./apps/next/public

ARG PORT=3002
ENV PORT=${PORT}
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
EXPOSE ${PORT}


# Add debugging commands
CMD ["bun", "apps/next/server.js"]
