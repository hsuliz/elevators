{ pkgs, lib, config, inputs, ... }:
{
  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    which

    pkgs.esbuild
    pkgs.nodePackages.prettier
  ];

  # https://devenv.sh/languages/
  languages = {
    nix.enable = true;

    go = {
      enable = true;
      package = pkgs.go_1_25;
    };

    typescript.enable = true;

    javascript = {
      enable = true;
      npm.enable = true;
    };
  };

  git-hooks.hooks = {
    #golangci-lint.enable = true;
    golines.enable = true;
    gotest.enable = true;
    nixpkgs-fmt.enable = true;
  };

  enterShell = ''
    go version
    echo 'go root' $GOROOT
    echo 'go path' $GOPATH
  '';

  scripts = {
    nix-format.exec = "nixpkgs-fmt .";

    go-install.exec = "go mod download";
    go-build.exec = "go build ./...";

    ts-build.exec = "tsc";
    ts-watch.exec = "esbuild web/main.ts --bundle --outfile=static/main.js --watch";
  };
}
