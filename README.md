# Slopify

> [!NOTE]  
> I created this project "For Science!". The goal was to find out how hard it is to completely automate low-effort content creation... Turns out it's really not that hard.

![image](https://github.com/notarock/slopify/assets/25652765/8ccb648e-1ed6-45da-a9e8-75d84137c696)

## Introduction

Slopify is an innovative tool designed to transform Reddit comment threads into
short, engaging videos. By leveraging natural language processing and video
generation technologies, Slopify offers users a unique way to consume and share
glimpses, jokes, stories, and discussions from Reddit in a visually appealing
format

Thanks you chatGPT, bery cool.

## What it can actually do

- **Create brainrot short-fort content videos :** Automatically converts Reddit comment threads into short videos with text and narration.
- **Customizable video background :** Offers customizable backgrounds to make each video unique using `--footage`
- **Easy upload :** Built-in sharing options to easily distribute your generated videos on platforms like YouTube.
- **Built-in title generation :** Prompt ChatGPT for SEO title and description to turn the [Dead Internet Theory](https://en.wikipedia.org/wiki/Dead_Internet_theory) into a real thing!

## Usage

``` sh
go run main.go reddit <permalink to comment / thread> --footage <path to background content folder>
```

## Project Setup

### Requirements
- Go 1.18
- ffmpeg
- terraform/opentofu
- A GCP project.
- OpenAPI key and some credits

#### 1. Clone repository

```sh
git clone https://github.com/notarock/slopify

```

#### 2. Configure GCP project

``` sh
gcloud auth login
gcloud config set project slopify
gcloud auth application-default login

```

#### 3. Create GCS storage bucket using provided terraform code

``` sh
cd terraform
terraform init
terraform apply
```

GCS storage is only used temporarily during the execution of the program. To use
the transcription API to obtain subtitles, the video passed as a parameter must
be in GCS...

The video is deleted in all cases using a defer, so it should not incur
additional costs related to storage usage.

#### 4. Activate GCP API

Activate these required APIS:
- Cloud Text-To-Speech API : https://console.cloud.google.com/apis/api/texttospeech.googleapis.com
- Cloud Video Intelligence https://console.cloud.google.com/apis/api/videointelligence.googleapis.com
- Youtube Data API v3

<img width="833" alt="image" src="https://github.com/notarock/slopify/assets/25652765/f589af56-927e-4874-a7b5-8888e48a114d">

#### 5. Add yourself as a developer in the project's users and create an OAuth application for YouTube.

<img width="947" alt="image" src="https://github.com/notarock/slopify/assets/25652765/6cd689d3-8cdf-43e2-9aa7-186125b13e87">

<img width="561" alt="image" src="https://github.com/notarock/slopify/assets/25652765/d984ac79-9b8c-45cd-85e5-0596c7f20de3">

<img width="895" alt="image" src="https://github.com/notarock/slopify/assets/25652765/5434c09e-0f18-4159-bab6-c2a1464f639e">

#### 6. Authorize YouTube through the OAuth app. The CLI will ask for the code found in the app's callback URL.

Download the JSON configuration file for the OAuth app and run the program,
which will initiate the authentication flow for YouTube. Then, copy the code
from the callback URL and paste it into the terminal.
