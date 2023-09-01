{ pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  buildInputs = with pkgs; [ xorg.libX11 xorg.libXcursor xorg.libXrandr xorg.libXinerama xorg.xinput xorg.libXi libGLU xorg.libXxf86vm ];
}

