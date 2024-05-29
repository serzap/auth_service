# Auth Service API Documentation

## Overview

The Auth Service API provides a set of methods to manage user authentication, including registration, login, email verification, and retrieval of user information. The API is implemented using gRPC.

## API Methods

### Register

Registers a new user.

#### Request

**RegisterRequest**

- `email` (string): The user's email address.
- `password` (string): The user's password.
- `username` (string): The user's username.
- `first_name` (string): The user's first name.
- `last_name` (string): The user's last name.

#### Response

**RegisterResponse**

- `user_id` (int64): The unique ID of the registered user.

### Login

Logs in a user and returns an authentication token.

#### Request

**LoginRequest**

- `email` (string): The user's email address.
- `password` (string): The user's password.

#### Response

**LoginResponse**

- `token` (string): The authentication token.

### VerifyEmail

Verifies a user's email address.

#### Request

**VerifyEmailRequest**

- `email` (string): The user's email address.
- `verification_code` (string): The verification code sent to the user's email.

#### Response

**VerifyEmailResponse**

- `success` (bool): Whether the email verification was successful.

### GetUserInfo

Retrieves user information based on the authentication token.

#### Request

**GetUserInfoRequest**

- `token` (string): The authentication token.

#### Response

**GetUserInfoResponse**

- `user_id` (int64): The unique ID of the user.
- `email` (string): The user's email address.
- `username` (string): The user's username.
- `first_name` (string): The user's first name.
- `last_name` (string): The user's last name.

## Running the Project

To run the project, follow these steps:

1. **Configure the settings**: configure authv1.yaml file containing parameters such as database connection details, RPC server addresses, and other settings.

2. **Set up the database**: Either run the database manually or use Docker. If running manually, ensure that your database server is running and accessible. If using Docker, run the following command: `docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=auth_service -e MYSQL_USER=user -e MYSQL_PASSWORD=password -p 3306:3306 -d mysql:latest`

3. **Start the server**: Launch the gRPC server for your service using the compiled binary or the run command.

4. **Testing the API**: You can test the API using tools like Postman or by running automated tests located in the `test` directory.
