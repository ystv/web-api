# TO DO

## Functionality
### Endpoints
- [ ] API for public website
  - [ ] VOD
    - [x] Videos
    - [x] Series
    - [ ] Playlists
    - [ ] Search
    - [x] Breadcrumbs
    - [x] Path to video
    - [x] Path to series
    - [ ] Thumbnail inheritance
  - [ ] Live
    - [x] Streams (currrently static)
  - [ ] Teams
    - [x] Get current team
    - [x] List teams
    - [ ] Get a team by year
    - [x] List all officers
- [ ] Internal (Secured)
  - [ ] Creator Studio
    - [ ] Videos
      - [x] Listing
      - [x] Getting with files
      - [x] Uploading
      - [ ] Updating
      - [ ] Deleting
    - [x] Series (very experimental)
    - [ ] Playlists
      - [x] List meta
      - [x] Create meta
      - [x] Update meta
      - [ ] Create video list
      - [ ] Update video list
    - [ ] Encoding
      - [ ] Presets
        - [x] Get
        - [x] Create (experimental)
        - [x] Update (experimental)
        - [ ] Delete
      - [ ] Formats
        - [x] Get
        - [ ] Create
        - [ ] Update
        - [ ] Delete
    - [x] Calendar
    - [x] Stats
  - [ ] People
    - [x] User by ID
    - [x] User by JWT
    - [x] Permissions
    - [ ] Create user
    - [ ] Create roles
    - [x] List all users
    - [x] List all users by role
  - [ ] Clapper
    - [ ] Events
      - [x] List by month
      - [ ] List by term
      - [ ] Signups
        - [x] Listing including positions
        - [x] Creating
        - [x] Updating
    - [ ] Postitions
      - [x] List
      - [x] Create
      - [x] Update
      - [ ] Delete
  - [ ] Encoder
    - [x] Video upload auth hook
  - [x] Stream auth (experimental)
  - [ ] Misc internal services

### Services
- [ ] Encode management
- [x] Mailer

### Connections
- [x] Postgres
- [x] RabbitMQ
- [x] web-auth intergration
- [ ] Authstack

## Documentation
- [ ] Install process
- [ ] Code structure
- [ ] Data routing
