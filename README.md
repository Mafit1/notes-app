# Notes API Server

Notes API server on Go for taking notes with PostgreSQL database and JWT authentication.

## Getting started

### Setup
1. Clone repository:
   ```
   git clone https://github.com/Mafit1/notes-app.git
   cd notes-app
   ```

2. Edit .env.example and rename it to .env

3. Install Docker if you haven't already

4. Run Docker compose:
   ```
   docker-compose up --build
   ```

### Metrics
You can get server metrics by open grafana on http://localhost:3000 and prometheus on http://localhost:9090

### API Endpoints

| Method | Endpoint   | Description          | Sample request body                                                |
| ------ | ---------- | -------------------- | ------------------------------------------------------------------ |
| POST   | /register  | Registration         | {"name": "user","email": "email@gmail.com","password": "123456789" |
| POST   | /login     | Login                | {"email": "email@gmail.com","password": "123456789"                |
| POST   | /refresh   | Refresh tokens       | -                                                                  |
| POST   | /logout    | Logout               | -                                                                  |
| POST   | /notes     | Create note          | {"title": "Shop list","content": "Bread, milk"}                    |
| GET    | /notes     | Get all notes        | -                                                                  |
| GET    | /notes/:id | Get note with id     | -                                                                  |
| PUT    | /notes/:id | Update note          | {"title": "Shop list","content": "Bread, milk, eggs"}              |
| DELETE | /notes/:id | Delete note          | -                                                                  |
