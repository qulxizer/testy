FROM node:22-alpine AS base

WORKDIR /app


FROM base AS deps

COPY package.json package-lock.json ./
RUN npm ci


FROM base AS dev

COPY --from=deps /app/node_modules ./node_modules
COPY package.json package-lock.json ./

ENV NODE_ENV=development

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--hostname", "0.0.0.0"]


FROM base AS builder

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=production

RUN npm run build


FROM base AS prod

ENV NODE_ENV=production

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

CMD ["node", "server.js"]
