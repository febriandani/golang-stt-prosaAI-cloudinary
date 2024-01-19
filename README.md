# Golang API for Speech-to-Text with ProsaAI and Cloudinary

This Golang API is designed to perform speech-to-text conversion using ProsaAI and allows users to upload audio files to Cloudinary for processing. It leverages the power of ProsaAI's speech-to-text capabilities and Cloudinary's file storage and processing features.

## Getting Started

Follow the steps below to set up and run the Golang API on your local machine.

### Prerequisites

1. Go installed on your machine. [Download and install Go](https://golang.org/dl/)

2. ProsaAI API key. [Sign up for ProsaAI](https://prosa.ai/) to get your API key.

3. Cloudinary account. [Sign up for Cloudinary](https://cloudinary.com/users/register/free) to obtain your Cloudinary API key, API secret, and Cloud name.

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/febriandani/golang-stt-prosaAI-cloudinary
    cd your-repo
    ```

2. Clone the repository:

    ```bash
    go mod tidy
    ```

3. Create a `config/app.yaml` file in the project root and add your ProsaAI and Cloudinary credentials:

    ```env
    CLOUDINARY:
        CLOUDINARY_URL : cloudinary://
        CLOUD_API_KEY: 418138124125412512321
        CLOUD_API_SECRET: Vb13fBX
        CLOUD_CLOUDNAME: ds14123
        URL_DOMAIN: https://api.cloudinary.com/v1_1/

    PROSA:
        API_KEY: eyJhbGciOiJSUzI1NiIsImtpZCI6Ik5XSTBNemRsTXprdE5tSmtNa
        URL: wss://s-api.prosa.ai/v2/speech/stt
    ```

4. Run the application:

    ```bash
    make run
    ```

The API should now be running at `http://localhost:8004`.


## Feedback

If you have any feedback, please reach out to us at febriandani176@gmail.com


## 🔗 Links
[![portfolio](https://img.shields.io/badge/my_portfolio-000?style=for-the-badge&logo=ko-fi&logoColor=white)](https://github.com/febriandani/)
[![linkedin](https://img.shields.io/badge/linkedin-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/mhmmdfebriandani//)

