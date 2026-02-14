#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://api:3000}"
DB_URL="${DB_URL:-postgresql://platform:platform_secret@db:5432/platform?sslmode=disable}"
LOG_FILE="${SMOKE_LOG_FILE:-/tmp/smoke_api.log}"
CHECKLIST_FILE="${SMOKE_CHECKLIST_FILE:-/tmp/endpoint_checklist.md}"
EXPECTED_FILE="${ENDPOINTS_EXPECTED_FILE:-/app/scripts/endpoints_expected.txt}"
ROUTES_HIT_FILE="$(mktemp)"

: > "$LOG_FILE"
: > "$CHECKLIST_FILE"

trim() {
  tr -d '[:space:]'
}

sql_scalar() {
  local sql="$1"
  psql "$DB_URL" -v ON_ERROR_STOP=1 -AtqX -c "$sql" | head -n1 | trim
}

request() {
  local method="$1"; shift
  local path="$1"; shift
  local body="${1:-}"
  local token="${2:-}"
  local route_template="${3:-$path}"

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
    echo "route: $method $route_template"
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

  echo "$method $route_template" >> "$ROUTES_HIT_FILE"
}

echo "[smoke] checking public endpoints" >&2
request GET "/" "" "" "/" >/dev/null
request GET "/api/docs/" "" "" "/api/docs/" >/dev/null
request GET "/api/docs/swagger.json" "" "" "/api/docs/swagger.json" >/dev/null

echo "[smoke] auth/register/login" >&2
REG='{"organizationName":"QA Org","inn":"1231231231","contactEmail":"qa-org@example.com","contactPhone":"+79990000001","adminName":"QA Admin","adminEmail":"qa-admin@example.com","adminPassword":"password123"}'
reg_resp=$(request POST "/auth/register" "$REG" "" "/auth/register")
TOKEN=$(echo "$reg_resp" | jq -r '.accessToken')

login_resp=$(request POST "/auth/login" '{"email":"qa-admin@example.com","password":"password123"}' "" "/auth/login")
TOKEN2=$(echo "$login_resp" | jq -r '.accessToken')
[[ -n "$TOKEN" && "$TOKEN" != "null" ]]
[[ -n "$TOKEN2" && "$TOKEN2" != "null" ]]

me_resp=$(request GET "/auth/me" "" "$TOKEN" "/auth/me")
ORG_ID=$(echo "$me_resp" | jq -r '.organizationId')
USER_ID=$(echo "$me_resp" | jq -r '.id')

request GET "/users" "" "$TOKEN" "/users" >/dev/null
request GET "/users/${USER_ID}" "" "$TOKEN" "/users/{id}" >/dev/null
request PATCH "/users/${USER_ID}" '{"name":"QA Admin Updated"}' "$TOKEN" "/users/{id}" >/dev/null
request GET "/organizations/${ORG_ID}" "" "$TOKEN" "/organizations/{id}" >/dev/null
request PATCH "/organizations/${ORG_ID}" '{"name":"QA Org Updated"}' "$TOKEN" "/organizations/{id}" >/dev/null

EXTRA_USER_ID=$(sql_scalar "INSERT INTO users(email,password_hash,name,role,organization_id) VALUES ('qa-user@example.com','x','QA User','MANAGER','$ORG_ID') RETURNING id;")
request DELETE "/users/${EXTRA_USER_ID}" "" "$TOKEN" "/users/{id}" >/dev/null

contractor_resp=$(request POST "/contractors" '{"name":"Contractor One","phone":"+79990000002","type":"INDIVIDUAL","regions":["Moscow"],"equipmentTypes":["Excavator"],"consentGiven":true}' "$TOKEN" "/contractors")
CONTRACTOR_ID=$(echo "$contractor_resp" | jq -r '.id')
request GET "/contractors" "" "$TOKEN" "/contractors" >/dev/null
request GET "/contractors/${CONTRACTOR_ID}" "" "$TOKEN" "/contractors/{id}" >/dev/null
request PATCH "/contractors/${CONTRACTOR_ID}" '{"isAvailable":true,"priceExpectations":"1000"}' "$TOKEN" "/contractors/{id}" >/dev/null
bl_resp=$(request POST "/contractors/${CONTRACTOR_ID}/blacklist" '{"reason":"test"}' "$TOKEN" "/contractors/{id}/blacklist")
BLACKLIST_ID=$(echo "$bl_resp" | jq -r '.id')
request GET "/contractors/${CONTRACTOR_ID}/history" "" "$TOKEN" "/contractors/{id}/history" >/dev/null

job_resp=$(request POST "/jobs" '{"title":"Excavator Work","region":"Moscow","equipmentType":"Excavator"}' "$TOKEN" "/jobs")
JOB_ID=$(echo "$job_resp" | jq -r '.id')
request GET "/jobs" "" "$TOKEN" "/jobs" >/dev/null
request GET "/jobs/${JOB_ID}" "" "$TOKEN" "/jobs/{id}" >/dev/null
request PATCH "/jobs/${JOB_ID}" '{"conditions":"updated"}' "$TOKEN" "/jobs/{id}" >/dev/null
request GET "/matching/jobs/${JOB_ID}/contractors" "" "$TOKEN" "/matching/jobs/{jobId}/contractors" >/dev/null
request POST "/jobs/${JOB_ID}/publish" "" "$TOKEN" "/jobs/{id}/publish" >/dev/null
request POST "/jobs/${JOB_ID}/complete" "" "$TOKEN" "/jobs/{id}/complete" >/dev/null
DUP_JOB_ID=$(request POST "/jobs/${JOB_ID}/duplicate" "" "$TOKEN" "/jobs/{id}/duplicate" | jq -r '.id')
request POST "/jobs/${DUP_JOB_ID}/cancel" "" "$TOKEN" "/jobs/{id}/cancel" >/dev/null

ASSIGNMENT_ID=$(request POST "/assignments" "{\"jobId\":\"${JOB_ID}\",\"contractorId\":\"${CONTRACTOR_ID}\"}" "$TOKEN" "/assignments" | jq -r '.id')
request GET "/assignments/job/${JOB_ID}" "" "$TOKEN" "/assignments/job/{jobId}" >/dev/null
request PATCH "/assignments/${ASSIGNMENT_ID}/status" '{"status":"CONFIRMED"}' "$TOKEN" "/assignments/{id}/status" >/dev/null

request POST "/ratings" "{\"jobId\":\"${JOB_ID}\",\"contractorId\":\"${CONTRACTOR_ID}\",\"score\":5,\"comment\":\"great\"}" "$TOKEN" "/ratings" >/dev/null
request GET "/ratings/contractor/${CONTRACTOR_ID}" "" "$TOKEN" "/ratings/contractor/{contractorId}" >/dev/null
request GET "/ratings/job/${JOB_ID}" "" "$TOKEN" "/ratings/job/{jobId}" >/dev/null

request GET "/compliance/blacklist" "" "$TOKEN" "/compliance/blacklist" >/dev/null
request DELETE "/compliance/blacklist/${BLACKLIST_ID}" "" "$TOKEN" "/compliance/blacklist/{id}" >/dev/null
request GET "/analytics/dashboard" "" "$TOKEN" "/analytics/dashboard" >/dev/null

DISPATCH_ID=$(sql_scalar "INSERT INTO job_dispatches(job_id,wave,status,\"limit\") VALUES ('$JOB_ID',1,'PENDING',10) RETURNING id;")
TOKEN_A=$(sql_scalar "INSERT INTO job_offers(job_id,dispatch_id,contractor_id,status) VALUES ('$JOB_ID','$DISPATCH_ID','$CONTRACTOR_ID','SENT') RETURNING token;")
TOKEN_D=$(sql_scalar "INSERT INTO job_offers(job_id,dispatch_id,contractor_id,status) VALUES ('$DUP_JOB_ID','$DISPATCH_ID','$CONTRACTOR_ID','SENT') RETURNING token;")

request GET "/responses/${TOKEN_A}" "" "" "/responses/{token}" >/dev/null
request POST "/responses/${TOKEN_A}/accept" "" "" "/responses/{token}/accept" >/dev/null
request GET "/responses/${TOKEN_D}" "" "" "/responses/{token}" >/dev/null
request POST "/responses/${TOKEN_D}/decline" "" "" "/responses/{token}/decline" >/dev/null

sort -u "$ROUTES_HIT_FILE" -o "$ROUTES_HIT_FILE"

{
  echo "# Endpoint Checklist"
  echo
} > "$CHECKLIST_FILE"

missing=0
while IFS= read -r endpoint; do
  [[ -z "$endpoint" ]] && continue
  if grep -Fxq "$endpoint" "$ROUTES_HIT_FILE"; then
    echo "- [x] $endpoint" >> "$CHECKLIST_FILE"
  else
    echo "- [ ] $endpoint" >> "$CHECKLIST_FILE"
    missing=1
  fi
done < "$EXPECTED_FILE"

echo "CHECKLIST_FILE=${CHECKLIST_FILE}"

echo "[smoke] endpoint checklist:"
cat "$CHECKLIST_FILE"
echo "SMOKE_LOG_FILE=${LOG_FILE}"

if [[ "$missing" -ne 0 ]]; then
  echo "[smoke] endpoint coverage check failed: some routes are not tested" >&2
  exit 1
fi

echo "[smoke] all endpoint checks passed"

rm -f "$ROUTES_HIT_FILE"
