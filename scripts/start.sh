#!/bin/bash

# Function to prompt user for using default .env files or creating new ones
prompt_env_files() {
  read -p "Do you want to use the default .env files? (y/n): " choice
  case "$choice" in
    y|Y )
      echo "Using default .env files..."
      cp configs/Docker.env.template configs/Docker.env
      cp frontend/.env.template frontend/.env
      ;;
    n|N )
      echo "Please create your .env files in 'configs/Docker.env' and 'frontend/.env', you can use configs/Docker.env.template and frontend/.env.template as templates. Do not forget to copy configs/Docker.env to configs/.env."
      exit 0
      ;;
    * )
      echo "Invalid choice. Please enter y or n."
      prompt_env_files
      ;;
  esac
}

# Prompt user for .env files
prompt_env_files

# Copy Docker.env to .env
cp configs/Docker.env configs/.env

echo "Setup complete."