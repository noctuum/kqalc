{
  description = "Full qalculate power in KRunner";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      forAllSystems = nixpkgs.lib.genAttrs [ "x86_64-linux" "aarch64-linux" ];
      pkgsFor = system: nixpkgs.legacyPackages.${system};
    in
    {
      packages = forAllSystems (system:
        let pkgs = pkgsFor system; in
        {
          default = pkgs.buildGoModule {
            pname = "kqalc";
            version = self.shortRev or self.dirtyShortRev or "dev";
            src = self;

            # Update: nix build 2>&1 | grep 'got:'
            vendorHash = "sha256-Ac63bZlBvCrhS7b8mk7aJdApI8UGtJxnZG35L37roGY=";

            env.CGO_ENABLED = 0;
            ldflags = [ "-s" "-w" ];

            nativeBuildInputs = [ pkgs.makeWrapper ];

            postInstall = ''
              install -Dm644 dist/org.kde.krunner1.kqalc.desktop \
                $out/share/krunner/dbusplugins/org.kde.krunner1.kqalc.desktop

              install -Dm644 dist/org.kde.krunner1.kqalc.service \
                $out/share/dbus-1/services/org.kde.krunner1.kqalc.service
              substituteInPlace $out/share/dbus-1/services/org.kde.krunner1.kqalc.service \
                --replace-warn "/usr/bin/kqalc" "$out/bin/kqalc"
            '';

            meta = with pkgs.lib; {
              homepage = "https://github.com/noctuum/kqalc";
              license = licenses.gpl2Only;
              platforms = [ "x86_64-linux" "aarch64-linux" ];
            };
          };
        }
      );
    };
}
