# Proof Of Work challenge-response protocol server

### Environment variables:  
SERVER_NETWORK_TYPE - type of network for server (default tcp);  
WISDOM_WORDS_FILE_PATH - path to wisdom words file;
USER_FILE_PATH - path to users key-value storage file;
SERVER_ADDRESS - address of server (default :8080);  

### Restrictions
Server could not be started without USER_FILE_PATH and WISDOM_WORDS_FILE_PATH env variables.   
By default, user file is expected to be in format key=value, each pair - new line.   
````
user=password
user1=password1
````
By default, wisdom words file is expected to be in format:
````
If you want to achieve greatness stop asking for permission. ~Anonymous
Things work out best for those who make the best of how things work out. ~John Wooden
````
