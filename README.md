# BIND API

This project is a Restful API for BIND DNS Servers made in GO. 

**THIS PROJECT IS CURRENTLY NOT UPDATING DNS, JUST STORING USERS**

## Setting up Project

We recommend to use the docker-compose of current project.
After installing **docker** and **docker-compose** , 
you need to put your external IP on dns image ``docker-compose.yml`` .
After configuring the file, you just need to run the following command
and the system will be on:

```
docker-compose up
``` 

## Services

There are four services running in this project: 

- Bind: The DNS server which exposes the default port 53
- Webmin: The admin web to manipulate the dns Server which is running on port 10000
- MariaDB: The database service to store users. It runs on port 8083
- BackEnd: the backend that store the users and whatever, which runns on port 8080

The Bind and Webmin runs on the same container.

##Backend API
The backend API is very simple. It create users with a hashed bcrypt password 
and is authenticated via JWT. It currently has four routes:

- POST /register: This route create a new user to use the DNS Server.
 It receive the following parameters on body:
```
{ 
 username: username of the new user to be created,
 password: password of the user that will be created,
 email:  email of the user that will be created
}
```
- POST /login: This route logins the current user. It receiver the following paramenters:

```
{ 
 username: username that will be logged,
 password: password of the user that will be logged,
}
```

it returns a json with username and jwt token with 5 minutes to expire on both answer and cookie: 
``` 
{ 
 username: username that was logged,
 token: JWT token that was generated. This has 5 minutes to expire
}
```

- POST /refresh: this route is only used to refresh the current JWT token.
this only need to be used to refresh the token when it's about to expire.
It returns the same type of answer as login:

``` 
{ 
 username: username that was logged,
 token: JWT token that was generated. This has 5 minutes to expire
}
```
- POST /domains: ** NOT WORKING ** this route is used to update a DNS server using NSUPDATE.
