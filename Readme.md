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


### shorten logic:
1. Receive long URL from user
2. we get a ID from a table name `url_id_generator` where we will have multiple rows with ranges of ID to chose from what to pick
 ex - row 1: start_id: 1, end_id: 1,000,000
        row 2: start_id: 1,000,001, end_id: 2,000,000
3. Pick a row from the table `url_id_generator` where the range is not exhausted
 ```PostgreSQL
    SELECT * FROM url_id_generator WHERE current_id < end_id order by RANDOM() LIMIT 1 FOR UPDATE;
    UPDATE url_id_generator SET current_id = current_id + 1 WHERE id = <picked_row_id>;
 ```
4. Use the `current_id` from the picked row to generate a short URL using base62 encoding
5. Store the mapping of long URL to short URL in the Database
6. Return the short URL to the user

### redirect logic:
1. Receive short URL from user
2. Check Redis Cache for the short URLs
3. If found in Cache, return the corresponding long URLs
4. If not found in Cache, query the Database for the short URL
5. If found in Database, store the mapping in Redis Cache for future requests
6. Return the long URL to the user

### Rate Limiting:
- Implemented using Redis to track the number of requests per user/IP
- Set a limit of 50 requests per minute per user/IP
- If the limit is exceeded, return a 429 Too Many Requests response
- Using Fixed Window Counter algorithm for rate limiting

