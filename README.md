# UpptimePrivateBackend

- Update `.env` file to add your PAT Token that you used for Upptime.
- Update your `.upptimerc.yml` to include following under `status-website:`:
```yaml
  apiBaseUrl: http://localhost:8080 #private
  userContentBaseUrl: http://localhost:8080/raw #private
  publish: true #private
```
- Change localhost:8080 with your domain/ip
- Enjoy
