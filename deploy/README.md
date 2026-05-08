# COMPX576 Deployment

This deployment targets the temporary `contabo_576` k3s server.

## Server order

1. Verify SSH key login.
2. Add UFW allow rules for 22, 6443, and Cloudflare-only 80/443.
3. Set UFW default deny incoming and allow outgoing.
4. Enable UFW.
5. Disable SSH password login.
6. Enable fail2ban.
7. Install k3s with Traefik disabled.
8. Install ingress-nginx.
9. Apply namespace, configmap, Postgres, Redis, and workloads.
10. Configure Cloudflare, cert-manager, issuer, and ingress last.

Do not test fail2ban by banning this local machine.

## Deploy scripts

Set:

```bash
export GHCR_OWNER=<owner>
export API_BASE_URL=https://<domain>/api
```

Run a script to deploy that component:

```bash
deploy/scripts/deploy-web-server.sh
deploy/scripts/deploy-agent-server.sh
deploy/scripts/deploy-agent.sh
deploy/scripts/deploy-frontend.sh
deploy/scripts/deploy-report.sh
```

Other operations are manual `kubectl` commands.
