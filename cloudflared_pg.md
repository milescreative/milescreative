# Cloudflare Tunnel + Postgres Access (Remote Network Guide)

This guide walks through how to securely connect to your self-hosted Postgres instance via a Cloudflare Tunnel — even when you're not on your local network.

---

## 🔧 First-Time Setup (One-Time Only)

### On Your Server

1. **Install cloudflared**:
   ```bash
   curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o /usr/local/bin/cloudflared
   chmod +x /usr/local/bin/cloudflared
   ```

2. **Login to Cloudflare**:
   ```bash
   cloudflared tunnel login
   ```

3. **Create the tunnel**:
   ```bash
   cloudflared tunnel create postgres-tunnel
   ```

4. **Route to your Postgres service** (running on the server at `127.0.0.1:5433`):
   ```bash
   cloudflared tunnel route dns postgres-tunnel pg.mlcr.us
   ```

5. **Create a config file at `~/.cloudflared/config.yml`**:
   ```yaml
   tunnel: postgres-tunnel
   credentials-file: /root/.cloudflared/<your-tunnel-id>.json

   ingress:
     - hostname: pg.mlcr.us
       service: tcp://127.0.0.1:5433
     - service: http_status:404
   ```

6. **Start the tunnel**:
   ```bash
   cloudflared tunnel run postgres-tunnel
   ```

---

### In Cloudflare Dashboard (First Time Only)

1. Go to **Access → Applications → Add an Application**
2. Choose **Self-hosted**
3. **Subdomain**: `pg.mlcr.us`
4. Leave path blank
5. Add a policy: Allow your email (e.g., `your@email.com` or `*@gmail.com`)
6. Save

---

## 💻 On Your Laptop

### One-Time Setup

1. **Install cloudflared**:
   ```bash
   brew install cloudflared
   ```

2. **Login to Cloudflare Access**:
   ```bash
   cloudflared access login https://pg.mlcr.us
   ```

---

## 🔁 Reconnecting Later (Remote Access Instructions)

Every time you need to connect from a non-local network:

1. **Start the tunnel proxy on your laptop**:
   ```bash
   cloudflared access tcp --hostname pg.mlcr.us --url localhost:5433
   ```

2. **Connect to Postgres** using a client (e.g., psql, TablePlus):

   ```
   Host: localhost
   Port: 5433
   User: postgres
   Password: <your_password>
   Database: postgres
   SSL: require
   ```

   Or with `psql`:
   ```bash
   psql "postgres://postgres:<your_password>@localhost:5433/postgres?sslmode=require"
   ```

---

## 🛠 Troubleshooting

- Make sure the **server is running**:
  ```bash
  cloudflared tunnel run postgres-tunnel
  ```

- If connecting from a new device/network, re-run:
  ```bash
  cloudflared access login https://pg.mlcr.us
  ```

- Confirm nothing else is using port 5433 on your laptop:
  ```bash
  lsof -i :5433
  ```
