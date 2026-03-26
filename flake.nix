{
  description = "A Emacs like compilation-mode for TUI";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.buildGoModule {
            pname = "gobuild";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-uFyrLQ5Q3sE2ioqdYxxfXgua8pG9ebIDcYKV3nVN568=";

            meta = with pkgs.lib; {
              description = "Emacs like compilation-mode for TUI";
              homepage = "https://github.com/xendak/gobuild";
              license = licenses.mit;
              mainProgram = "gobuild";
            };
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
            ];

            shellHook = ''
              export GOPATH="$PWD/.go_cache"
              export GOMODCACHE="$GOPATH/pkg/mod"
              export PATH="$GOPATH/bin:$PATH"

              mkdir -p "$GOMODCACHE"
            '';
          };
        }
      );
    };
}
