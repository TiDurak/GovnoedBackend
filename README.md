<h1 align="center">
    <img src="https://user-images.githubusercontent.com/82606298/170775456-475ffa71-9cf9-4584-9723-b3917ae0aecc.svg" alt="DebilBot" border="0" height="30px"> 
    DebilBot
</h1>


# GovnoedBackend

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)  

![Discord](https://img.shields.io/discord/439778332121497612?style=flat-square&logo=discord&logoColor=white&label=Discord%20Server)
![Website](https://img.shields.io/website?url=https%3A%2F%2Fgovnoed.de%2F&up_message=Visit&down_message=Not%20working&style=flat-square&label=govnoed.de)
![GitHub](https://img.shields.io/website?url=https%3A%2F%2Fgithub.com%2FTiDurak%2Fdebilbot&up_message=tidurak%2Fdebilbot&up_color=black&down_message=removed&style=flat-square&logo=github&label=GitHub)


## Description

This project is a refactored backend for the [govnoed.de](https://govnoed.de/) website.

The API uses an external SQLite database, communicates with the website through HTTP API endpoints, and generates promo codes for [debilbot](https://github.com/tidurak/debilbot).


## Project Status

The project is in its final stage. Future bug fixes and other minor fixes may still be made.


## Project Structure

The backend is organized into three independent microservices:

- **GovnoedPromo**: Handles promo code generation and management for the debilbot integration
- **GovnoedUserItems**: Manages user items and related data
- **GovnoedWeb**: Provides the web interface and serves the website frontend


## Architecture

This project employs a **microservices architecture** to maintain stability and ease of development. Each service is a standalone application that can be developed, tested, and deployed independently.

### Communication

The services communicate with each other through **HTTP API endpoints**. This decoupled design offers several advantages:

- **Scalability**: Services can be scaled independently based on demand
- **Maintainability**: Each service has a focused responsibility and a contained codebase
- **Resilience**: Failure of one service doesn't necessarily bring down the entire system
- **Development Flexibility**: Different services can use different databases and technologies if needed
- **Easy Testing**: Individual services can be tested in isolation

### Technology Stack

All services are built in **Go** for:
- High performance with minimal resource consumption
- Simple, readable code compared to other systems languages
- Fast compilation and deployment
- Built-in concurrency support for handling multiple requests


## FAQ

Q: Why was the backend rewritten?  
A: As of 23 August 2026, the current [govnoed.de](https://govnoed.de/get_key) backend is written in PHP. PHP was the best choice for my website at the time. However, the project quickly became difficult to scale and turned into an unreadable mess.


Q: Why Go instead of Rust, Django, or Node.js?  
A: The goal was to build a lightweight API because my VPS has very limited resources. The choice ultimately came down to Rust and Go. I chose Go because it keeps the code much simpler with only a small loss in performance. I also had no prior experience with Rust.
