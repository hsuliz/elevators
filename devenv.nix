{ pkgs, lib, config, inputs, ... }:
{
  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    which
  ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
    package = pkgs.go_1_25;
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
    buildApp.exec = ''
      go build ./...
    '';
  };
}
