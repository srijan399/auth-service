# JWT Auth + RBAC — Conceptual Flow (Go)

## 1. Schema

- `users` — id, email, password_hash, created_at
- `roles` — id, name (e.g. 'user', 'admin')
- `user_roles` — user_id, role_id (many-to-many join)

Use a join table even with one role per user today — cheap insurance against a future migration.

## 2. Registration flow

1. Client sends email + password.
2. Hash password with bcrypt — never store plaintext or a reversible hash.
3. Insert into `users`.
4. Insert a row into `user_roles` linking the new user to the default role (`'user'`). This is a deliberate app decision — nobody self-registers as admin.

## 3. Login flow

1. Client sends email + password.
2. Look up user by email.
3. Compare submitted password against stored hash (bcrypt compare, not equality on hashes).
4. If valid: look up the user's roles, build a claims payload (subject/user id, roles, expiry, issued-at).
5. Sign the claims with a secret (HMAC) or private key (asymmetric) to produce the JWT.
6. Return the token to the client — either in the JSON response body, or as an httpOnly cookie. Pick one architecture deliberately; don't mix without reason.

## 4. Storing the secret

- Never hardcode the signing secret/key in source.
- Load from environment variable or secrets manager.
- Same secret/key pair must be available wherever you verify tokens (could be multiple services if asymmetric signing is used).

## 5. Sending the token on future requests

- If body-delivered: client attaches `Authorization: Bearer <token>` manually on every request.
- If cookie-delivered: browser attaches it automatically.
- Either way — JWT auth is stateless. The server has no memory of "you logged in earlier." Every protected request re-proves identity via the token itself.

## 6. Authentication middleware (runs on every protected route)

Two responsibilities, both required:

1. **Extract** the token from the request (header or cookie).
2. **Verify**: recompute the signature using the known secret/key, compare against the token's signature, and check expiry. If this fails — reject immediately (401/403), don't proceed.
3. On success: pull claims out of the verified token and attach them to the request context (user id, roles) so downstream handlers/middleware can use them without re-parsing.

Verification answers exactly one question: _did the server actually issue this, unmodified?_ Nothing about permissions yet.

## 7. Authorization middleware (separate from authentication)

- Reads roles already placed in context by the authentication step.
- Checks whether the required role for this route is present.
- If missing → 403 Forbidden.
- If present → call next handler.

Keep this as its own composable middleware (e.g. parameterized by required role), not folded into authentication. Authentication asks "who are you," authorization asks "are you allowed here" — different questions, different middleware.

## 8. Key design decisions to make consciously

| Decision       | Options                                                               | Tradeoff                                                                                  |
| -------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| Token delivery | JSON body vs httpOnly cookie                                          | XSS exposure vs CSRF exposure                                                             |
| Signing method | Symmetric (HS256) vs Asymmetric (RS256/ES256)                         | Single secret shared everywhere vs private key signs, public key verifies across services |
| Roles location | Baked into JWT claims vs looked up from DB per request                | Fast but stale until token expires, vs always fresh but costs a query                     |
| Token lifetime | Short-lived access token (+ refresh token) vs long-lived single token | Security/responsiveness to role changes vs convenience                                    |

## 9. The full request lifecycle, end to end

```
Register → hash password → insert user → assign default role
Login    → verify password → fetch roles → sign JWT → return token
Request  → client attaches token → Authenticate middleware verifies signature/expiry
         → claims (roles) placed in context → Authorize middleware checks role
         → handler runs only if both pass
(repeat verification on every single request — nothing is cached server-side)
```

## 10. Things to research next, once base flow works

- Algorithm confusion attacks (always pin/check the expected signing method during verify)
- Refresh token rotation and revocation strategy
- What happens when a role changes mid-token-lifetime (acceptable staleness window?)
- Logout semantics for stateless JWTs (there's no server session to destroy — consider short expiry + client-side token discard, or a blocklist if you need hard revocation)
