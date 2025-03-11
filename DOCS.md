# Chirpy API Docs

## Chirp Resource

```json
{
  "id": "624cdb2f-9bdd-4556-a64c-00986e3d1de7",
  "created_at": "2025-03-11T02:55:11.114722Z",
  "updated_at": "2025-03-11T02:55:11.114722Z",
  "email": "agustin@example.com",
  "is_chirpy_red": false
}
```

### GET /api/chirps

Returns an array of chirps

### GET /api/chirps/{chirpID}

Returns the Chirp with id `chirpID` if exists.

### POST /api/chirps

Receives a `body` and creates a new chirp with that body
(user needs to be authenticated).

Example

```json
{
  "body": "This is the text of the new chirp."
}
```

### DELETE /api/chirps/{chirpID}

Deletes the Chirp with id `chirpID` if exists.
