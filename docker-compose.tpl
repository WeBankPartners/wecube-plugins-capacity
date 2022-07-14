version: '2'
services:
  capacity:
    image: wecube-plugins-capacity:{{version}}
    container_name: wecube-plugins-capacity-{{version}}
    restart: always
    volumes:
      - /etc/localtime:/etc/localtime
      - {{path}}/capacity/logs:/app/capacity/logs
    ports:
      - "{{port}}:9096"
    environment:
      - CAPACITY_LOG_LEVEL={{log_level}}
      - CAPACITY_MYSQL_HOST={{db_server}}
      - CAPACITY_MYSQL_PORT={{db_port}}
      - CAPACITY_MYSQL_USER={{db_user}}
      - CAPACITY_MYSQL_PWD={{db_pass}}
      - CAPACITY_MYSQL_SCHEMA={{db_database}}
      - GATEWAY_URL={{monitor_url}}