# builder stage
FROM node:lts-alpine AS builder
WORKDIR /app

# copy package files and install dependencies with yarn
COPY webui/package.json webui/yarn.lock* ./
RUN yarn install

# now copy the rest of the webui code and build it
COPY webui/ .
RUN yarn run build


# final stage, nginx serves the stuff
FROM nginx:stable-alpine

# copy the built frontend from the builder
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80