# Nginx Setup Guide

## Requirements

- Ubuntu server with Nginx installed
- Domain DNS A record pointing to your server IP
- Certbot installed (`apt install certbot python3-certbot-nginx`)

---

## 1. Check Nginx Version

```bash
nginx -v
```

---

## 2. Create a Site Config

The Nginx config is already in this repo at `nginx/sites-available/hospital-a.api.co.th`. Copy it to the server:

```bash
sudo cp ~/agnos/nginx/sites-available/hospital-a.api.co.th /etc/nginx/sites-available/hospital-a.api.co.th
```

Or if cloning the repo on the server directly:

```bash
sudo cp /path/to/repo/nginx/sites-available/hospital-a.api.co.th /etc/nginx/sites-available/
```

The config file looks like this (before SSL — Certbot will modify it automatically):

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name hospital-a.api.co.th;

    location / {
        proxy_pass http://localhost:8082;
    }
}
```

> Change `proxy_pass` port to match your `APP_PORT`.

---

## 3. Enable the Site (Symlink)

```bash
sudo ln -s /etc/nginx/sites-available/hospital-a.api.co.th /etc/nginx/sites-enabled/
```

Verify the link was created:

```bash
ls /etc/nginx/sites-enabled
```

---

## 4. Test & Reload Nginx

```bash
sudo nginx -t
sudo systemctl restart nginx
```

---

## 5. Issue SSL Certificate via Certbot

```bash
sudo certbot --nginx -d hospital-a.api.co.th
```

Certbot will automatically modify the site config to add HTTPS. The final config will look like this:

```nginx
server {
    server_name hospital-a.api.co.th;

    location / {
        proxy_pass http://localhost:8082;
    }

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/hospital-a.api.co.th/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/hospital-a.api.co.th/privkey.pem;
}

server {
    listen 80;
    server_name hospital-a.api.co.th;
    return 301 https://$host$request_uri;
}
```

---

## 6. Test & Reload Again

```bash
sudo nginx -t
sudo systemctl restart nginx
```
