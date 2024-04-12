# Web-Messenger-Backend

Welcome to the Web-Messenger-Backend repository! This project serves as the backend for a web-based messenger application, facilitating real-time communications via robust API services designed to manage user interactions and data securely.

## Features
- User authentication (login, registration, logout).
- Profile management (retrieve and update user profile).
- Real-time messaging.
- Contacts and chat management.
- Multimedia messages and file transfers support.

## Technologies Used
- **Go**: For efficient backend services.
- **Docker**: For deployment and environment consistency.
- **WebSocket**: For real-time communication challenges.

## API Endpoints
- **Authentication**:
  - `/login`: Authenticate and return a session token.
  - `/register`: Handle new user registrations.
  - `/logout`: End user sessions.
- **Profile Management**:
  - `/user/profile`: Access and update user profiles.
- **Messaging**:
  - `/message/send`: Endpoint to send new messages.
  - `/message/receive`: Endpoint to receive messages.
- **Contacts & Chats**:
  - `/contacts`: Manage contact lists.
  - `/chats`: Retrieve chat histories.

## WebSocket Challenges

Implementing WebSockets was crucial for maintaining real-time user chats in our Web-Messenger-Backend. This necessitated upgrading from a traditional HTTP connection to a WebSocket connection to support continuous, two-way interactions essential for live messaging.

### Key Challenges:
- **Connection Upgrade**: Transitioning from HTTP to WebSocket involves managing upgrade headers, which can introduce compatibility issues with proxies and firewalls.
- **Message Integrity and Ordering**: Ensuring that messages are delivered reliably and in the correct order is critical, especially to prevent message loss or duplication.
- **Connection Management**: Managing each user's connection state is complex, particularly in environments with load balancers or when handling reconnects after connection drops.
- **Resource Management**: WebSockets maintain persistent connections, increasing server resource consumption. Efficient resource management is essential to maintain performance as user numbers grow.

These challenges required careful architectural decisions to ensure robust, scalable, and secure real-time communication for users.

## Getting Started
```bash
git clone https://github.com/aref81/Web-Messenger-Backend.git
cd Web-Messenger-Backend
docker-compose up --build
