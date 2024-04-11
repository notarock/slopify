{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
  buildInputs = [
    (google-cloud-sdk.withExtraComponents [ google-cloud-sdk.components.gke-gcloud-auth-plugin ])
    opentofu
    ffmpeg
    yt-dlp
    ollama
    imagemagickBig
    fontconfig
  ];


  shellHook = ''
    echo Performing gcloud auth login
    echo gcloud auth login
    echo gcloud config set project slopify
    echo gcloud auth application-default login
    export FONTCONFIG_FILE=${pkgs.fontconfig.out}/etc/fonts/fonts.conf
    export FONTCONFIG_PATH=${pkgs.fontconfig.out}/etc/fonts/
    export PKG_CONFIG_PATH=${imagemagickBig.dev}/lib/pkgconfig
  '';
}
