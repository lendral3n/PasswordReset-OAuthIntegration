# PasswordReset-OAuthIntegration

This project is an email notification service that provides user authentication features, including password reset and email verification. This project also supports OAuth authentication with Google and Facebook.

## 🔮 Features

- 👤 **Authentication and User Management**
  - User Registration
  - User Login
  - Get User Details
  - Update User Account
  - Update User Password
  - Delete User
  - Forgot Password
  - Reset Password via Email Link
  - Reset Password via Email Code
  - Email Verification via Email Link
  - Email Verification via Email Code
  - OAuth with Google
  - OAuth with Facebook

## Endpoint List

| Tag    | Endpoint                         |
| ------ | -------------------------------- |
| 👤User | `POST /login`                    |
| 👤User | `POST /users`                    |
| 👤User | `GET /users`                     |
| 👤User | `PUT /users`                     |
| 👤User | `DELETE /users`                  |
| 👤User | `PUT /change-password`           |
| 👤User | `POST /forgot-password`          |
| 👤User | `PATCH /reset-password`          |
| 👤User | `POST /verification`             |
| 👤User | `PATCH /verification`            |
| 👤User | `POST /request-code-password`    |
| 👤User | `PATCH /reset-password-code`     |
| 👤User | `POST /request-code-verify`      |
| 👤User | `PATCH /verification-email`      |
| 👤User | `GET /oauth-google`              |
| 👤User | `GET /api/sessions/oauth/google` |
| 👤User | `GET /oauth-facebook`            |
| 👤User | `GET /id/oauth/callback/`        |

## 🛠️ Technology Stack

- **Golang**: The programming language used to develop the backend of the application.
- **Echo**: A fast and easy-to-use web framework for Go, used to build web applications and APIs.
- **GORM**: An Object-Relational Mapping (ORM) library for Golang.
- **PostgreSQL**: A relational database management system.
- **Redis**: An in-memory data store that is super fast, used as a database, cache, and message broker.
- **Amazon S3**: An object storage service that offers scalability, data protection, and performance.
- **Amazon RDS**: A service that makes it easy to set up, operate, and scale a relational database in the cloud.
- **JWT**: JSON Web Token for securing data exchange.
- **AWS**: Amazon Web Services, a cloud platform that offers a variety of IT infrastructure services.
- **Docker**: A platform that allows developers to build, package, and distribute applications with ease.
- **Mailtrap**: A fake SMTP server for development teams to test, view and share emails sent from the development and staging environments.
- **Google Cloud**: A suite of cloud computing services that runs on the same infrastructure that Google uses internally for its end-user products.
- **Facebook Developer**: Tools and services to help you bring your ideas to life, reach an audience and monetize your app.

## Requirements

To run the project successfully, you need to set up various credentials and configurations for different services and components. Below are the requirements along with instructions on how to obtain the credentials:

### Database Configuration
```
DBUSER => The username for your database.
DBPASS => The password for your database.
DBHOST => The host of your database.
DBPORT => The port of your database.
DBNAME => The name of your database.
```

These credentials are required to connect to your PostgreSQL database. You can obtain them from your database administrator or by setting up a new database instance on your preferred hosting service.

### JWT Configuration
```
JWTSECRET => The secret key for signing JWTs.
```

You can generate a JWT Secret of your choice to secure your JWT tokens. Make sure it is a long, randomly generated string.

### Redis Configuration
```
RDSURL => The URL for your Redis instance.
```

For Redis caching, you need to set up a Redis instance. Obtain the Redis URL from your Redis hosting service.

### AWS S3 Configuration
```
AWSKEY => Your AWS Key ID.
AWSSECRET => Your AWS Secret Key.
AWSREGION => Your AWS S3 Region.
```

If you're using AWS S3 for storing multimedia assets, you need to create an AWS IAM user with S3 access. Obtain the AWS Key ID, AWS Secret Key, and AWS S3 Region from your AWS IAM user dashboard.

### Email Configuration
```
EMAILFROM => The email address that will be used as the sender in the emails.
SMTPHOST => The host for your SMTP server.
SMTPPORT => The port for your SMTP server.
SMTPUSER => The user for your SMTP server.
SMTPPASS => The password for your SMTP server.
```
For sending transactional emails, you can use services like Mailtrap or your own SMTP server. Obtain the SMTP credentials from your email service provider.

### Password Reset Configuration
```
PASSWDURL => The URL that will be used for password resets.
```

Define the URL where users can reset their passwords. This URL should point to your password reset endpoint in the project.

### Google OAuth Configuration


```
CLIENTID => Your Google Client ID for OAuth.
CLIENTSECRET => Your Google Client Secret for OAuth.
GOOGLEURL => Your Google Callback URL for OAuth.
SCOPES => Your Google Scopes for OAuth.
```

To integrate Google OAuth authentication, you need to create a project on the Google Developer Console and obtain OAuth credentials (Client ID and Client Secret). Configure the OAuth consent screen with the appropriate scopes and redirect URIs.

### Facebook OAuth Configuration
```
CLIENTIDFB => Your Facebook Client ID for OAuth.
CLIENTSECRETFB => Your Facebook Client Secret for OAuth.
FBURL => Your Facebook Callback URL for OAuth.
SCOPESFB => Your Facebook Scopes for OAuth.
```
Similar to Google OAuth, to integrate Facebook OAuth authentication, you need to create a project on the Facebook Developer Console and obtain OAuth credentials (Client ID and Client Secret). Configure the OAuth consent screen with the appropriate scopes and redirect URIs.

## 🧰 Installation
Follow these steps to install and set up the KosKita API:
1. **Clone the repository:**

   ```bash
   git clone https://github.com/lendral3n/PasswordReset-OAuthIntegration.git
   
2. **Move to Cloned Repository Folder**

    ```bash
    cd PasswordReset-OAuthIntegration
    
3. **Update dependecies**
    
    ```bash
    go mod tidy
    ```

4. **Copy `local.env.example` to `local.env`**

    ```bash
    cp local.env.example local.env
    ```

5. **Configure your `local.env`**
6. **Run PasswordReset-OAuthIntegration** 
7. 
    ```bash
    go run main.go
    ```
## 🤖 Author

- **Lendra Syaputra**
  - [Github](https://github.com/lendral3n)