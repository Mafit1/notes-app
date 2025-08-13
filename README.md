# Notes API Server

Notes API server on Go for taking notes with PostgreSQL database.

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

### API Endpoints

| Method | Endpoint   | Description          | Sample request body                                   |
| ------ | ---------- | -------------------- | ----------------------------------------------------- |
| POST   | /notes     | Create note          | {"title": "Shop list","content": "Bread, milk"}       |
| GET    | /notes     | Get all notes        | -                                                     |
| GET    | /notes/:id | Get note with id     | -                                                     |
| PUT    | /notes/:id | Update note          | {"title": "Shop list","content": "Bread, milk, eggs"} |
| DELETE | /notes/:id | Delete note          | -                                                     |
