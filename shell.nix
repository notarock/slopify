{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    (google-cloud-sdk.withExtraComponents
      [ google-cloud-sdk.components.gke-gcloud-auth-plugin ])

      ffmpeg

  ];

  shellHook = ''
    echo Performing gcloud auth login
    echo gcloud auth login
    echo gcloud config set project memes-traduit
    echo gcloud auth application-default login
  '';
}
