#!/usr/bin/env sh
# Smoke test of the key-value EXAMPLE (replace with calls to your service).
# Start `make run` (and optionally `make run2`) first.
# Usage: ./etc/demo.sh [base-url]     default: http://localhost:8080
B=${1:-http://localhost:8080}
J='Content-Type: application/json'
set -x

# --- GetAllBooks ---------------------------------------------------------
curl -si "$B/v1/AllBooks"
curl -si "$B/v1/AllBooks?author=sara"
curl -si "$B/v1/AllBooks?year=2018"
curl -si "$B/v1/AllBooks?limit=2"
curl -si "$B/v1/AllBooks?year=1900"
curl -si "$B/v1/AllBooks?consistency=linearizable"

# --- GetBook -------------------------------------------------------------
curl -si "$B/v1/AllBooks/Dummy%201?consistency=sequential"
curl -si "$B/v1/AllBooks/Dummy%201?consistency=linearizable"
curl -si "$B/v1/AllBooks/Unknown" | head -1  # 404

# --- AddBook -------------------------------------------------------------
curl -si -X POST "$B/v1/AllBooks" -H "$J" -d '{"Booktitle":"New Book","author":"New Author","year":2025,"Copies":4}'
curl -si -X POST "$B/v1/AllBooks" -H "$J" -d '{"Booktitle":"Dummy 1","author":"Author","year":2025,"Copies":4}' | head -1  # 409
curl -si -X POST "$B/v1/AllBooks" -H "$J" -d '{"Booktitle":"New Book","author":"Author","year":2025,"Copies":0}' | head -1  # 400

# --- AddCopies -----------------------------------------------------------
curl -si -X PATCH "$B/v1/AllBooks/Dummy%201" -H "$J" -d '{"copies":3}'
curl -si -X PATCH "$B/v1/AllBooks/Dummy%201" -H "$J" -d '{"copies":-3}' | head -1  # 400
curl -si -X PATCH "$B/v1/AllBooks/Unknown" -H "$J" -d '{"copies":3}' | head -1  # 404

# --- RemoveBook ----------------------------------------------------------
curl -si -X DELETE "$B/v1/AllBooks/Dummy%202" | head -1  # 204
curl -si -X DELETE "$B/v1/AllBooks/Unknown" | head -1  # 404
curl -si -X DELETE "$B/v1/AllBooks/Dummy%201" | head -1  # 409

# --- BorrowBook ----------------------------------------------------------
curl -si -X POST "$B/v1/AllBooks/Dummy%201/borrow" -H "$J" -d '{"userId":"user2"}'
curl -si -X POST "$B/v1/AllBooks/Dummy%201/borrow" -H "$J" -d '{"userId":""}' | head -1  # 400
curl -si -X POST "$B/v1/AllBooks/Unknown/borrow" -H "$J" -d '{"userId":"user2"}' | head -1  # 404
curl -si -X POST "$B/v1/AllBooks/Dummy%203/borrow" -H "$J" -d '{"userId":"user2"}' | head -1  # 409
curl -si -X POST "$B/v1/AllBooks/Dummy%201/borrow" -H "$J" -d '{"userId":"user1"}' | head -1  # 409


# --- ReturnBook ----------------------------------------------------------
curl -si -X POST "$B/v1/AllBooks/Dummy%201/return" -H "$J" -d '{"userId":"user1"}'
curl -si -X POST "$B/v1/AllBooks/Dummy%201/return" -H "$J" -d '{"userId":""}' | head -1  # 400
curl -si -X POST "$B/v1/AllBooks/Unknown/return" -H "$J" -d '{"userId":"user1"}' | head -1  # 404
curl -si -X POST "$B/v1/AllBooks/Dummy%201/return" -H "$J" -d '{"userId":"user2"}' | head -1  # 409


# --- BuyBook -------------------------------------------------------------
curl -si -X POST "$B/v1/AllBooks/Dummy%201/buy"
curl -si -X POST "$B/v1/AllBooks/Unknown/buy" | head -1  # 404
curl -si -X POST "$B/v1/AllBooks/Dummy%203/buy" | head -1  # 409

# --- docs -----------------------------------------------------------------
curl -s  "$B/openapi.yaml" | head -3
curl -si "$B/docs"   | head -1                                        # Swagger UI (open in a browser)
curl -si "$B/editor" | head -1                                        # Swagger Editor (open in a browser)
