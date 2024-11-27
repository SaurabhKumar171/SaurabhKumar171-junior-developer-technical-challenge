
# **Rick and Morty Character Search Engine**

This project is a full-stack application that allows users to search for characters from the popular *Rick and Morty* TV show. It leverages cloud infrastructure, a backend API built with Go, and a frontend search interface built with React and TypeScript.

---

## **Project Overview**

### **1. Cloud Infrastructure**
- **GCP Compute Engine**: Used to create a virtual machine with MongoDB Community Edition.
- **Terraform**: Used to automate infrastructure deployment, including setting up compute resources, block storage, and firewall rules.
- **Backup and Restore**: MongoDB data is backed up and can be restored from block storage when needed.

### **2. Backend API**
- **Go-based REST API**: Handles fetching character details from MongoDB and exposes a search functionality to search characters by name.
- **MongoDB Integration**: The database is populated with character data from the *Rick and Morty* API and stored persistently.

### **3. Frontend**
- **React & TypeScript**: A modern and responsive interface for searching characters and displaying their details and associated episodes.
- **TailwindCSS**: Used for styling the frontend with a minimal and modern design.

---

## **Directory Structure**

```plaintext
.
├── /terraform       # Contains Terraform configurations for setting up GCP infrastructure
├── /backend         # Golang backend that interacts with MongoDB
├── /frontend        # React + TypeScript frontend
└── README.md        # This file
```

---

## **Setup Instructions**

### **1. Prerequisites**

- **Google Cloud Platform**:
  - Active Google Cloud project with billing enabled.
  - Service account with Compute Engine and Storage permissions.
  - Credentials JSON file (e.g., `application_default_credentials.json`).
  
- **Required Tools**:
  - [Terraform](https://developer.hashicorp.com/terraform/docs) for provisioning infrastructure.
  - [Go](https://golang.org/doc/install) for running the backend API.
  - [Node.js](https://nodejs.org/) for running the frontend.

---

### **2. Setting Up Cloud Infrastructure (Terraform)**

#### Navigate to the Terraform Directory:
```bash
cd terraform
```

#### Configure Variables:
Open `variables.tf` and set the following variables:
- **project_id**: Your GCP project ID.

Open `main.tf` and set the following variables:
- **gcp_credentials_path**: Path to your Google Cloud credentials file.
- **source_ip**: Set your IP or `0.0.0.0/0` for unrestricted access.

#### Initialize and Apply Terraform Configuration:
```bash
terraform init
terraform plan -var="project_id=<your_project_id>"
terraform apply -var="project_id=<your_project_id>"
```

#### Outcome:
- A VM with MongoDB Community Edition deployed.
- A database `rickAndMorty` with a `characters` collection populated with character data from the *Rick and Morty* API.

#### Destroy Infrastructure (Optional):
To clean up:
```bash
terraform destroy -var="project_id=<your_project_id>"
```

---

### **3. Backend Setup**

#### Navigate to the Backend Directory:
```bash
cd backend
```

#### Open the .env file in the backend directory and set the MONGO_URL to the VM's IP address:
MONGO_URL=mongodb://<VM_IP>:27017

#### Install Dependencies:
```bash
go mod tidy
```

#### Run the Backend API:
```bash
go run main.go
```

The backend will now be running and can be accessed at `http://localhost:8000/api`.

#### API Endpoints:

- **Search Character by Name**:  
  `GET /api/characters/search?name=<character_name>`  
  Fetches details for the character with the specified name from MongoDB.

  **Example Request**:
  ```http
  http://localhost:8000/api/characters/search?name=Rick
  ```

- **Get Paginated List of Characters**:  
  `GET /api/getCharacters?page=<page>&limit=<limit>`  
  Fetches a paginated list of characters with the given page number and limit per page.

  **Example Request**:
  ```http
  http://localhost:8000/api/getCharacters?page=1&limit=6
  ```

---

### **4. Frontend Setup**

#### Navigate to the Frontend Directory:
```bash
cd frontend
```

#### Install Dependencies:
```bash
npm install
```

#### Start the Frontend Development Server:
```bash
npm start
```

The frontend will be available at `http://localhost:3000`.

---

## **Architecture Overview**

### **Infrastructure**
- **Cloud**: GCP Compute Engine VM with MongoDB.
- **Storage**: Persistent block storage for MongoDB data.
- **Networking**: Firewall rules to secure the MongoDB instance.

### **Backend**
- **Language**: Go
- **Database**: MongoDB for storing character data.
- **API**: RESTful API that supports searching characters.

### **Frontend**
- **Framework**: React with TypeScript for strong typing and scalability.
- **Styling**: TailwindCSS for a modern, responsive UI.

---

## **Features**

- **Character Search**: Search for characters by name using the `/api/characters/search` endpoint.
- **Paginated Results**: Fetch paginated character lists using the `/api/getCharacters` endpoint.
- **MongoDB Integration**: Data is stored and queried from MongoDB.
- **Data Persistence**: Block storage for backup and restoration.
- **Frontend Interface**: Built with React, TypeScript, and styled with TailwindCSS.
- **Error Handling**: Proper error states in the frontend, including loading and error messages.

---