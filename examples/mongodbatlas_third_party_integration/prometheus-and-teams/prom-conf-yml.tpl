global:
  scrape_interval: 15s
scrape_configs:

- job_name: "${job_name}"
  scrape_interval: 10s  
  metrics_path: /metrics
  scheme : https
  basic_auth:
    username: prom_user_${group_id}
    password: ${password}
  http_sd_configs:
    - url: https://cloud.mongodb.com/prometheus/v1.0/groups/${group_id}/discovery
      refresh_interval: 60s
      basic_auth:
        username: prom_user_${group_id}
        password: ${password}

