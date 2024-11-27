terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }

  required_version = ">= 1.0.0"
}

provider "google" {
  credentials = file(var.gcp_credentials_path)
  project     = var.project_id
  region      = var.region
  zone        = var.zone
}

# Create a Persistent Disk
resource "google_compute_disk" "mongodb_disk" {
  name  = "mongodb-vm-disk-3"
  type  = "pd-standard"
  size  = var.disk_size
  zone  = var.zone
}

# Create a Compute Instance for MongoDB
resource "google_compute_instance" "mongodb_vm" {
  name         = "mongodb-vm-test-3"
  machine_type = var.machine_type

  # Boot disk with OS
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts" # Base OS
    }
  }

  # Attach the persistent disk
  attached_disk {
    source = google_compute_disk.mongodb_disk.id
  }

  network_interface {
    network = "default"
    access_config {}
  }

  # Install and start MongoDB on VM creation
    metadata_startup_script = <<-EOT
    #!/bin/bash
    LOG_FILE="/var/log/mongodb_setup.log"

    echo "Starting MongoDB installation and configuration..." | tee -a $LOG_FILE
    sudo apt-get update >> $LOG_FILE 2>&1

    echo "Installing required dependencies..." | tee -a $LOG_FILE
    sudo apt-get install -y gnupg curl python3 python3-pip >> $LOG_FILE 2>&1

    echo "Importing MongoDB GPG key..." | tee -a $LOG_FILE
    curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | \
        sudo gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg --dearmor >> $LOG_FILE 2>&1

    echo "Adding MongoDB repository..." | tee -a $LOG_FILE
    echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/8.0 multiverse" | \
        sudo tee /etc/apt/sources.list.d/mongodb-org-8.0.list >> $LOG_FILE 2>&1

    echo "Updating package list..." | tee -a $LOG_FILE
    sudo apt-get update >> $LOG_FILE 2>&1

    echo "Installing MongoDB..." | tee -a $LOG_FILE
    sudo apt-get install -y mongodb-org >> $LOG_FILE 2>&1

    echo "Starting MongoDB service..." | tee -a $LOG_FILE
    sudo systemctl start mongod >> $LOG_FILE 2>&1

    echo "Enabling MongoDB service on boot..." | tee -a $LOG_FILE
    sudo systemctl enable mongod >> $LOG_FILE 2>&1

    echo "Binding MongoDB to all IPs..." | tee -a $LOG_FILE
    sudo sed -i 's/^\( *bindIp: \).*/\10.0.0.0/' /etc/mongod.conf >> $LOG_FILE 2>&1

    #echo "Enabling authorization in MongoDB..." | tee -a $LOG_FILE
    #sudo sed -i '/^#security:/s/^#security:/security:\n  authorization: enabled/' /etc/mongod.conf

    echo "Restarting MongoDB service..." | tee -a $LOG_FILE
    sudo systemctl restart mongod >> $LOG_FILE 2>&1

    echo "Installing Python requests library..." | tee -a $LOG_FILE
    sudo pip3 install requests pymongo >> $LOG_FILE 2>&1

    echo "Downloading and storing Rick and Morty characters..." | tee -a $LOG_FILE
    python3 <<'PYTHON_SCRIPT'
    import requests
    from pymongo import MongoClient

    # MongoDB configuration
    mongo_client = MongoClient("mongodb://localhost:27017")
    db = mongo_client["rickAndMorty"]
    collection = db["characters"]

    # Fetch characters from API
    def fetch_characters():
        url = "https://rickandmortyapi.com/api/character"
        all_characters = []
        while url:
            print(f"Fetching: {url}")
            response = requests.get(url)
            if response.status_code == 200:
                data = response.json()
                all_characters.extend(data["results"])
                url = data["info"]["next"]
            else:
                print(f"Failed to fetch: {url}")
                break
        return all_characters

    # Store characters in MongoDB
    characters = fetch_characters()
    if characters:
        collection.insert_many(characters)
        print(f"Inserted {len(characters)} characters into MongoDB.")
    PYTHON_SCRIPT

    echo "Setup complete!" | tee -a $LOG_FILE
    EOT


  tags = ["mongodb"]
}

# Firewall Rule to Allow MongoDB Access
resource "google_compute_firewall" "allow_mongo_new_3" {
  name    = "allow-mongo-new-3"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["27017"] # MongoDB default port
  }

  source_ranges = [var.source_ip]
}


