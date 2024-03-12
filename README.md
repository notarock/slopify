# Slopify

![image](https://github.com/notarock/slopify/assets/25652765/8ccb648e-1ed6-45da-a9e8-75d84137c696)

## Introduction
Slopify est un outil innovant conçu pour transformer les fils de commentaires
Reddit en vidéos courtes et engageantes. En exploitant les technologies de
traitement du langage naturel et de génération de vidéo, Slopify offre aux
utilisateurs une manière unique de consommer et partager des aperçus, blagues,
histoires et discussions de Reddit sous une forme visuellement attrayante.

Merci ChatGPT pour ce beau résumé

## Fonctionnalités
- **Création Automatique de Vidéos :** Convertit automatiquement les fils de commentaires Reddit en vidéos courtes avec texte et narration.
- **Visuels Personnalisables :** Propose des arrière-plans personnalisables pour rendre chaque vidéo unique. `--footage`
<!-- - **Options de Voix Off :** Comprend une variété d'options de voix off pour correspondre au ton et au contexte du contenu. -->
- **Partage Facile :** Options de partage intégrées pour distribuer facilement vos vidéos générées sur les plateformes tel que Youtube

## Installation

### Prérequis
- Go 1.18
- ffmpeg
- terraform/opentofu
- un project GCP
- Clef d'api OpenAI ainsi que quelques crédits'

### Étapes
1. Clonez le dépôt :
```sh
git clone https://github.com/notarock/slopify

```

2. Configurer le project GCP
``` sh
gcloud auth login
gcloud config set project slopify
gcloud auth application-default login

```

3. Créer le bucket de storage GCS.

``` sh
cd terraform
terraform init
terraform apply
```

Le storage GCS est seulement utilisé de manière temporaire
durant l'Execution du programme. Pour utiliser l'api de transcription afin
d'obtenir les sous-titres, il faut que le vidéo passé en param soit dans GCS...

Le vidéo se delete dans tous les cas à l'aide d'un `defer` donc ca ne devrait
pas causé de frais supplémentaire relié à l'utilisation de storage.

4. Activer les API Google Cloud pour le projets

Il faut activer ces API:
- Cloud Text-To-Speech API : https://console.cloud.google.com/apis/api/texttospeech.googleapis.com
- Cloud Video Intelligence https://console.cloud.google.com/apis/api/videointelligence.googleapis.com
- Youtube Data API v3

<img width="833" alt="image" src="https://github.com/notarock/slopify/assets/25652765/f589af56-927e-4874-a7b5-8888e48a114d">

5. S'ajouter en tant que développeur dans les users du project et créer une application oauth pour Youtube

<img width="947" alt="image" src="https://github.com/notarock/slopify/assets/25652765/6cd689d3-8cdf-43e2-9aa7-186125b13e87">

<img width="561" alt="image" src="https://github.com/notarock/slopify/assets/25652765/d984ac79-9b8c-45cd-85e5-0596c7f20de3">

<img width="895" alt="image" src="https://github.com/notarock/slopify/assets/25652765/5434c09e-0f18-4159-bab6-c2a1464f639e">

6. Authorizer youtube par l'app Oauth. Le cli va demander le code qui va se retrouver dans l'url de callback de l'app

Il faut télécharger le fichier de configuration json pour l'app oauth et lancer le programme, ca va
initier le flow de authentification pour youtube.  ensuite copier le code dans l'URL callback et le
coller dans le terminal
