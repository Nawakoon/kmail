## Roadmap

1. Server
 - health check
 - inbox
    * need to read total mail instead of the length of queried mail in page
 - send email
 - read email

2. Client
 - inbox
    * need to make it not require user private in command, have a choice to type as password in shell
 - send email
 - read email

3. Infrastructure
 - setup local database by postgres image docker ci
 - setup prod database by Aurora AWS ci
 - deploy server ci
 - host client image ci