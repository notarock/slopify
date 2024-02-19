{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    (google-cloud-sdk.withExtraComponents
      [ google-cloud-sdk.components.gke-gcloud-auth-plugin ])
      terraform
      ffmpeg
      yt-dlp
  ];

  shellHook = ''
    echo Performing gcloud auth login
    echo gcloud auth login
    echo gcloud config set project slopify
    echo gcloud auth application-default login
  '';
}
