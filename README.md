# Insta_API
A Api made using Go language without any 3rd party libraries, to clone the back end functionalities of Instagram. 

## Features

- Create User - (POST) /users
- Get a user using id - (GET) /users/:id
- Create a Post -  (POST) /posts
- Get all Posts in DB using Pagination - (GET) /posts?page=2
- Get a post using id - (GET) /posts/:id
- List all posts of a user - /posts/user/:userid


## Extra Features

- Pagination
- Password encryption before storing in DB



## Routes

### Create User - (POST) /users

Request Body 
```
{
    "name": "Arnab",
    "email": "Arnab@gmail.com",
    "password": "34nnnsk123k"
}
```

Success Response - The auto generated user id comes as response
```
{
    "1633803186081694000"
}
```



### Get a User using ID - (GET) /users/:id

Success Response - The user object comes as a response
```
{
    "name": "Arnab",
    "email": "Arnab@gmail.com",
    "ID": "1633803186081694000",
    "Password": "18a132fc2339074ba9beb552c1e3bf4944e1947fc83b109844dd7005627f9e53ebc5f1e0aa269f"
}
```


### Create a post - (POST) /posts

Request Body 
```
{
    "caption":"Nice image!",
    "imageurl":"imageURL",
    "userid":"1633803024127613400"
}
```

Success Response - Post Id is returned and timestamp is auto generated
```
{
    "6161df3cc145c2cec877026c"
}
```


### Get all posts in DB  - (GET) /posts?page=2

Success Response 
```
{[
    {
        "caption": "name of the coaster 20",
        "imageUrl": "",
        "ID": "1633799942177085800",
        "UserID": "12345",
        "timestamp": "0001-01-01T00:00:00Z"
    },
    {
        "caption": "",
        "imageUrl": "",
        "ID": "1633800895545488600",
        "UserID": "",
        "timestamp": "0001-01-01T00:00:00Z"
    },
    {
        "caption": "",
        "imageUrl": "",
        "ID": "1633803024127613400",
        "UserID": "",
        "timestamp": "0001-01-01T00:00:00Z"
    },
    {
        "caption": "",
        "imageUrl": "",
        "ID": "1633803843679088300",
        "UserID": "",
        "timestamp": "0001-01-01T00:00:00Z"
    }
]
```


### Get a post using ID - (GET) /posts/:id

Success Response - The post object comes as a response
```
{
    {
    "caption": "Nice image 3!",
    "imageUrl": "imageURL",
    "ID": "1633804585499709300",
    "UserID": "1633803024127613400",
    "timestamp": "0001-01-01T00:00:00Z"
}
```


### Get all posts of a user - (GET) /posts/users/:id


```
[
    {
        "caption": "Nice image!",
        "imageUrl": "imageURL",
        "ID": "1633804092363660300",
        "UserID": "1633803024127613400",
        "timestamp": "0001-01-01T00:00:00Z"
    },
    {
        "caption": "Nice image 2!",
        "imageUrl": "imageURL",
        "ID": "1633804509819333700",
        "UserID": "1633803024127613400",
        "timestamp": "0001-01-01T00:00:00Z"
    },
    {
        "caption": "Nice image 3!",
        "imageUrl": "imageURL",
        "ID": "1633804585499709300",
        "UserID": "1633803024127613400",
        "timestamp": "0001-01-01T00:00:00Z"
    }
]
```



## Tech Stack
 - Language - Gp
 - Database - Mongo DB



