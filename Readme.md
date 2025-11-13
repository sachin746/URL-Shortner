## URL-Shorter Service

#### Technologies Used: GoLang, GORM, PostgreSQL, Redis

 #### UseCase:
- long URL to short URL conversion
- short URL to long URL redirection

#### Features:
- Rate Limiting: Limits the number of requests a user can make in a given time period
- Caching: Uses Redis to cache frequently accessed URLs for faster retrieval
- Analytics: Tracks the number of clicks on each short URLs
- Expiration: Allows users to set an expiration date for their short URLs
- User Authentication: Users can create accounts and manage their URLs 
- Custom Aliases: Users can create custom short URLs instead of random accounts

#### TPS And Scale:
- per second requests: 1000
- total URLs to be handled: `1000*60*60*24*30*365` = 31,536,000,000 (31.5 billion URLs per year)
- total Reads/Writes: 90% Reads, 10% Writes
- Read Requests: 28,382,400,000 per year
- Write Requests: 3,153,600,000 per year

#### Database Schema:
- Url Table:
  - id (Primary Key)
  - long_url (Text)
  - short_url (Varchar)
  - user_id (Foreign Key to User Table)
  - created_at (Timestamp)
  - expires_at (Timestamp, Nullable)
  - click_count (Integer)
