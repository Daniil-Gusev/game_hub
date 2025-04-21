#!/bin/bash

if [ "$#" -ne 5 ]; then
    echo "Usage: $0 <platform: linux|darwin> <app_name> <version> <binary_path> <output_dir>"
    exit 1
fi

platform=$1
app_name=$2
app_version=$3
binary_path=$4
output_dir=$5

templates_dir="templates"
run_sh_template="$templates_dir/run.sh"
desktop_template="$templates_dir/app.desktop"
plist_template="$templates_dir/info.plist"

if [ ! -f "$run_sh_template" ] || [ ! -f "$desktop_template" ] || [ ! -f "$plist_template" ]; then
    echo "Error: Not all templates found in $templates_dir/"
    exit 1
fi

binary_name=$(basename "$binary_path")

mkdir -p "$output_dir"

render_template() {
    template_file=$1
    sed \
        -e "s/\$AppName/$app_name/g" \
        -e "s/\$AppVersion/$app_version/g" \
        -e "s/\$BinaryName/$binary_name/g" \
        "$template_file"
}

if [ "$platform" = "linux" ]; then
    desktop_file="$output_dir/$app_name.desktop"
    render_template "$desktop_template" > "$desktop_file"
    chmod +x "$desktop_file"
    echo "Created .desktop file: $desktop_file"

elif [ "$platform" = "darwin" ]; then
    app_bundle="$output_dir/$app_name.app"
    mkdir -p "$app_bundle/Contents/MacOS"
    mkdir -p "$app_bundle/Contents/Resources"    
    cp "$binary_path" "$app_bundle/Contents/Resources/$binary_name"
    run_sh_file="$app_bundle/Contents/MacOS/run.sh"
    render_template "$run_sh_template" > "$run_sh_file"
    chmod +x "$run_sh_file"
    plist_file="$app_bundle/Contents/Info.plist"
    render_template "$plist_template" > "$plist_file"
    echo "Created .app bundle: $app_bundle"
else
    echo "Unknown platform: $platform. Only linux or darwin are supported."
    exit 1
fi