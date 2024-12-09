# **Book Review Website** 💻

## ⭐️ Overview

This repository contains the backend logic for a media review application for books.

It uses:

```text
🔹 gin-gonic framework for the router
🔹 air for live loading
🔹 golang-jwt for tokans/cookies
🔹 MongoDB for the database storage
```

## ⚙️ Prerequisites

In order to run this application you will need to have the following:

```text
🔸 MongoDB account
🔸 VS code installed
🔸 Golang installed
🔸 air-verse/air live loader installed
```

## 🐾 Step One

Change your directory to where you wish to run this script and store the cloned repository:

```bash
cd <filename>
```

## 🐾 Step Two

Clone the repository from github and then move into the new directory.

```bash
git clone https://github.com/laurawarren88/spa_backend.git
cd spa_backend
```

## 🐾 Step Three

Take a look 👀 around the file 📂 structure and see what is happening with VS code.

```bash
code .
```

You will need to add a .env 🤫 file into the file tree to store the necessary information to run the script.

I'll walk you through it:

```bash
touch .env
vim .env
```

In the file you need to include your information ℹ️ into the following variables:

```text
ENV=development
DEV_ALLOWED_ORIGIN=http://localhost:<port number>
PROD_ALLOWED_ORIGIN=http://<domain>:<port number>

PORT=8080
MONGODB_URI=<your mongodb uri>
SECRET_KEY=<your secret key>
ADMIN_PASSWORD=<your admin password>
```

For the ENV variable you can use development or production. This will determine which port the server will run on, you can set these in the next variables. These are your frontend ports for either development or production. You can use the same port number for both. What ever you use for the port number will be the port number you will need to use in the frontend.

For the port number you can use any port you like but if you are running this with the frontend you will need to change the port number in the frontend file as well - /utils/config.js.

After the '=' sign for DATABASE_URL input the connection for your MongoDB, it will look something like this:

```text
mongodb+srv://<username>:<password>@cluster0.ib6l0.mongodb.net/<cluster_name>?retryWrites=true&w=majority&appName=Cluster0
```

For the secret, and admin_password variable input anything you like.

## 🐾 Step Four

Ensure the repository builds successfully, MongoDB is connected and the server is running, by running the following:

```bash
air
```

## 🐾 Step Five

Now the backend is running you can use Postman to test the API.

Input the following into the Postman URL bar:

```text
http://localhost:8080/api
```

This should load the home page. From here you can visit other routes by adding /books or /reviews. You can also test the sign up and login routes.

There is an admin account already set up with the following credentials:

```text
username: admin
Email:    "admin@admin.com",
password: <what you set in the .env file>
```

## 🐾 Step Six

In order to view the frontend of the application you will need to clone the frontend repository and run the application.

This can be found here:

```text
https://github.com/laurawarren88/spa_frontend
```

Follow the readme on how to run the application.
