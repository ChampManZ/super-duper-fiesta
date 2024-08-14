# UrMessage - Take-Home Assignment | Small Social Media

**This project repository is dedicated for take home assignment and self-education while developing only**

Welcome to UrMessage, a streamlined social media messaging application developed for a take-home assignment and educational purposes. The repository name has been randomly generated.

UrMessage functions as a minimalist platform where users can post messages and comments, drawing inspiration from Twitter and Threads.

## Features

UrMessage offers the following features for users and administrators:

- **User Account Management**: Create and manage user accounts, serving as the core identifier within the application.
- **Authentication System**: Secure login system with session token generation.
- **Public Feed**: Publicly accessible feed where unregistered users can view posts and registered users can contribute content.
- **User Profile Management**: Personal profile page for updating user information.
- **Admin Tools**: Admin-exclusive UI for managing database migrations.
- **Commenting System**: Interactive commenting functionality on individual posts.
- **State Management**: Enhanced user experience with loading indicators during data fetching.

## Project Overview

This project integrates both frontend and backend within a single repository. Feedback on the following aspects is highly valued:

- **Best Practices**: Guidance on organizing folder structures for projects that include both UI and server components within the same repository.
- **Security Enhancements**: Suggestions for improving security in line with best practices.
- **Performance Optimization**: Insights into optimizing application performance, particularly in data handling and UI responsiveness.
- **Code Readability vs. Performance**: Advice on balancing code readability with performance optimization.
- **Coding Style**: Recommendations for maintaining a consistent and readable coding style, especially in collaborative environments.

## Setup

To set up and run UrMessage locally:
 ```bash
docker-compose build
docker-compose up -d
 ```

## Run Test
 ```bash
go test .\tests\
 ```
 or you can in Docker container using
```bash
docker exec -it super-duper-fiesta-server-1 go test .\tests\
```

## Versions
- node v21.7.3
- npm 9.5.1
- go version go1.22.6 windows/amd64 (test locally on Windows)
- jwt-go v5

## Swagger UI
http://localhost:1323/swagger/index.html#/

# Discussion & Room for Improvement
- **Responsive Design**: Improve the UI for better responsiveness across devices.
- **Performance Optimization**: Enhance performance in tasks like retrieving database migration files, especially as the dataset grows.
- **Secret Management**: Better handling of secrets and tokens to ensure security.
- **State Management**: Consider using Redux for state management in larger applications.
- **UI Testing**: Implement comprehensive UI testing, including unit, component, integration, and end-to-end tests.

# This README got help using GPT for grammar checking and wording that's more appropriate
