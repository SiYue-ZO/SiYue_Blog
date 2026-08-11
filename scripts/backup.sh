#!/bin/sh
set -eu
umask 077

mkdir -p /backups

sql_target=""
temporary=""
cleanup() {
  [ -z "$sql_target" ] || rm -f "$sql_target"
  [ -z "$temporary" ] || rm -f "$temporary"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

while true; do
  timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
  sql_target="/tmp/siyue_blog_${timestamp}.sql"
  target="/backups/siyue_blog_${timestamp}.sql.gz"
  temporary="${target}.tmp"
  pg_dump --no-owner --no-privileges > "$sql_target"
  gzip -9 < "$sql_target" > "$temporary"
  mv "$temporary" "$target"
  rm -f "$sql_target"
  sql_target=""
  temporary=""
  find /backups -type f -name 'siyue_blog_*.sql.gz' -mtime "+${BACKUP_RETENTION_DAYS:-14}" -delete
  sleep "${BACKUP_INTERVAL_SECONDS:-86400}"
done
