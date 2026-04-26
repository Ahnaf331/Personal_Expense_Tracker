#!/bin/bash
BASE="http://localhost:8080/api"

echo "=== Register ==="
REGISTER=$(curl -s -X POST $BASE/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}')
echo $REGISTER | python3 -m json.tool 2>/dev/null || echo $REGISTER

TOKEN=$(echo $REGISTER | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo ""
echo "Token: $TOKEN"
AUTH="Authorization: Bearer $TOKEN"

echo ""
echo "=== Login ==="
curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}' | python3 -m json.tool 2>/dev/null

echo ""
echo "=== Get Categories ==="
curl -s $BASE/categories -H "$AUTH" | python3 -m json.tool 2>/dev/null

echo ""
echo "=== Create Expense ==="
EXPENSE=$(curl -s -X POST $BASE/expenses \
  -H "Content-Type: application/json" \
  -H "$AUTH" \
  -d '{"category_id":1,"amount":45.50,"description":"Groceries","date":"2026-04-26"}')
echo $EXPENSE | python3 -m json.tool 2>/dev/null || echo $EXPENSE

echo ""
echo "=== Get All Expenses ==="
curl -s $BASE/expenses -H "$AUTH" | python3 -m json.tool 2>/dev/null

echo ""
echo "=== Monthly Analytics ==="
curl -s "$BASE/analytics/monthly?year=2026&month=4" -H "$AUTH" | python3 -m json.tool 2>/dev/null

echo ""
echo "=== Category Breakdown ==="
curl -s "$BASE/analytics/categories" -H "$AUTH" | python3 -m json.tool 2>/dev/null
