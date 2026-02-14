#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://api:3000}"
DB_URL="${DB_URL:-postgresql://platform:platform_secret@db:5432/platform?sslmode=disable}"
LOG_FILE="${SMOKE_LOG_FILE:-/tmp/smoke_api.log}"

: > "$LOG_FILE"

trim() {
  tr -d '[:space:]'
}

request() {
  local method="$1"; shift
  local path="$1"; shift
  local body="${1:-}"
  local token="${2:-}"

  local url="${BASE_URL}${path}"
  local body_file
  body_file="$(mktemp)"

  local code
  if [[ -n "$body" && -n "$token" ]]; then
    code=$(curl -sS -o "$body_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d "$body")
  elif [[ -n "$body" ]]; then
    code=$(curl -sS -o "$body_file" -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" -d "$body")
  elif [[ -n "$token" ]]; then
    code=$(curl -sS -o "$body_file" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $token")
  else
    code=$(curl -sS -o "$body_file" -w "%{http_code}" -X "$method" "$url")
  fi

  {
    echo "=== $method $path"
    echo "status: $code"
    echo "body:"
    cat "$body_file"
    echo
  } | tee -a "$LOG_FILE" >&2

  cat "$body_file"
  rm -f "$body_file"

  if [[ "$code" -lt 200 || "$code" -ge 300 ]]; then
    echo "request failed: $method $path (status=$code)" >&2
    exit 1
  fi
}

echo "[smoke] checking public endpoints" >&2
request GET "/" >/dev/null
request GET "/api/docs/" >/dev/null
request GET "/api/docs/swagger.json" >/dev/null

echo "[smoke] auth/register/login" >&2
REG='{"organizationName":"QA Org","inn":"1231231231","contactEmail":"qa-org@example.com","contactPhone":"+79990000001","adminName":"QA Admin","adminEmail":"qa-admin@example.com","adminPassword":"password123"}'
reg_resp=$(request POST "/auth/register" "$REG")
TOKEN=$(echo "$reg_resp" | jq -r '.accessToken')

login_resp=$(request POST "/auth/login" '{"email":"qa-admin@example.com","password":"password123"}')
TOKEN2=$(echo "$login_resp" | jq -r '.accessToken')
[[ -n "$TOKEN" && "$TOKEN" != "null" ]]
[[ -n "$TOKEN2" && "$TOKEN2" != "null" ]]

me_resp=$(request GET "/auth/me" "" "$TOKEN")
ORG_ID=$(echo "$me_resp" | jq -r '.organizationId')
USER_ID=$(echo "$me_resp" | jq -r '.id')

request GET "/users" "" "$TOKEN" >/dev/null
request GET "/users/${USER_ID}" "" "$TOKEN" >/dev/null
request PATCH "/users/${USER_ID}" '{"name":"QA Admin Updated"}' "$TOKEN" >/dev/null
request GET "/organizations/${ORG_ID}" "" "$TOKEN" >/dev/null
request PATCH "/organizations/${ORG_ID}" '{"name":"QA Org Updated"}' "$TOKEN" >/dev/null

EXTRA_USER_ID=$(psql "$DB_URL" -Atc "INSERT INTO users(email,password_hash,name,role,organization_id) VALUES ('qa-user@example.com','x','QA User','MANAGER','$ORG_ID') RETURNING id;" | trim)
request DELETE "/users/${EXTRA_USER_ID}" "" "$TOKEN" >/dev/null

contractor_resp=$(request POST "/contractors" '{"name":"Contractor One","phone":"+79990000002","type":"INDIVIDUAL","regions":["Moscow"],"equipmentTypes":["Excavator"],"consentGiven":true}' "$TOKEN")
CONTRACTOR_ID=$(echo "$contractor_resp" | jq -r '.id')
request GET "/contractors" "" "$TOKEN" >/dev/null
request GET "/contractors/${CONTRACTOR_ID}" "" "$TOKEN" >/dev/null
request PATCH "/contractors/${CONTRACTOR_ID}" '{"isAvailable":true,"priceExpectations":"1000"}' "$TOKEN" >/dev/null
bl_resp=$(request POST "/contractors/${CONTRACTOR_ID}/blacklist" '{"reason":"test"}' "$TOKEN")
BLACKLIST_ID=$(echo "$bl_resp" | jq -r '.id')
request GET "/contractors/${CONTRACTOR_ID}/history" "" "$TOKEN" >/dev/null

job_resp=$(request POST "/jobs" '{"title":"Excavator Work","region":"Moscow","equipmentType":"Excavator"}' "$TOKEN")
JOB_ID=$(echo "$job_resp" | jq -r '.id')
request GET "/jobs" "" "$TOKEN" >/dev/null
request GET "/jobs/${JOB_ID}" "" "$TOKEN" >/dev/null
request PATCH "/jobs/${JOB_ID}" '{"conditions":"updated"}' "$TOKEN" >/dev/null
request GET "/matching/jobs/${JOB_ID}/contractors" "" "$TOKEN" >/dev/null
request POST "/jobs/${JOB_ID}/publish" "" "$TOKEN" >/dev/null
request POST "/jobs/${JOB_ID}/complete" "" "$TOKEN" >/dev/null
DUP_JOB_ID=$(request POST "/jobs/${JOB_ID}/duplicate" "" "$TOKEN" | jq -r '.id')
request POST "/jobs/${DUP_JOB_ID}/cancel" "" "$TOKEN" >/dev/null

ASSIGNMENT_ID=$(request POST "/assignments" "{\"jobId\":\"${JOB_ID}\",\"contractorId\":\"${CONTRACTOR_ID}\"}" "$TOKEN" | jq -r '.id')
request GET "/assignments/job/${JOB_ID}" "" "$TOKEN" >/dev/null
request PATCH "/assignments/${ASSIGNMENT_ID}/status" '{"status":"CONFIRMED"}' "$TOKEN" >/dev/null

request POST "/ratings" "{\"jobId\":\"${JOB_ID}\",\"contractorId\":\"${CONTRACTOR_ID}\",\"score\":5,\"comment\":\"great\"}" "$TOKEN" >/dev/null
request GET "/ratings/contractor/${CONTRACTOR_ID}" "" "$TOKEN" >/dev/null
request GET "/ratings/job/${JOB_ID}" "" "$TOKEN" >/dev/null

request GET "/compliance/blacklist" "" "$TOKEN" >/dev/null
request DELETE "/compliance/blacklist/${BLACKLIST_ID}" "" "$TOKEN" >/dev/null
request GET "/analytics/dashboard" "" "$TOKEN" >/dev/null

DISPATCH_ID=$(psql "$DB_URL" -Atc "INSERT INTO job_dispatches(job_id,wave,status,\"limit\") VALUES ('$JOB_ID',1,'PENDING',10) RETURNING id;" | trim)
TOKEN_A=$(psql "$DB_URL" -Atc "INSERT INTO job_offers(job_id,dispatch_id,contractor_id,status) VALUES ('$JOB_ID','$DISPATCH_ID','$CONTRACTOR_ID','SENT') RETURNING token;" | trim)
TOKEN_D=$(psql "$DB_URL" -Atc "INSERT INTO job_offers(job_id,dispatch_id,contractor_id,status) VALUES ('$DUP_JOB_ID','$DISPATCH_ID','$CONTRACTOR_ID','SENT') RETURNING token;" | trim)

request GET "/responses/${TOKEN_A}" >/dev/null
request POST "/responses/${TOKEN_A}/accept" >/dev/null
request GET "/responses/${TOKEN_D}" >/dev/null
request POST "/responses/${TOKEN_D}/decline" >/dev/null

echo "[smoke] all endpoint checks passed"
echo "SMOKE_LOG_FILE=${LOG_FILE}"
