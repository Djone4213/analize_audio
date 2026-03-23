# nginx/Dockerfile

# ---------- React build ----------
FROM node:20-trixie AS react-build

WORKDIR /app

COPY audio-client/package*.json ./
RUN npm install

COPY audio-client/ ./
RUN npm run build


# ---------- Nginx ----------
FROM debian:trixie-slim

RUN apt-get update && \
    apt-get install -y nginx ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=react-build /app/build /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
