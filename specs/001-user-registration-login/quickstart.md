# Quickstart: Local Validation

## Prerequisites

- Docker Desktop is running.
- `docker compose` is available.
- The workspace is the repository root.

## Start local services

```bash
docker compose up -d
```

## Validate the authentication flow

1. Open the web app locally.
2. Use the landing page to navigate to Create Account.
3. Submit valid registration data.
4. Confirm the app redirects to the login page.
5. Sign in using the newly created credentials.
6. Verify the authenticated home screen appears with the user’s first name and a logout action.
7. Trigger logout and confirm the signed-out session-ended state is shown.

## Expected outcome

The auth MVP should complete the full register → login → logout lifecycle without introducing any financial functionality.
