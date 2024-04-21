# Kitty-Sesh 🐱🚬

Kitty-sesh is an Kitty Terminal Session Manager, Go-based application that provides a graphical user interface (GUI) to manage and interact with Kitty terminal sessions. Users can save, view, rename, delete, and launch stored sessions directly from the application.

![image](https://github.com/Raghav-rv28/kitty-sesh/assets/62635473/70ae0a80-85b9-427b-9444-950cf7eafe0e)

## Installation

### Requirements

- Go (v1.16 or higher)
- Git
- Kitty Terminal

### Steps to install

- use the command below for installation in Linux.

```sh
bash -c "$(curl -sL https://raw.githubusercontent.com/Raghav-rv28/kitty-sesh/main/install.sh)"
```

## Usage

use the command `kitty-sesh` to start the application in your terminal.

- **Navigation**: Use Arrow keys to traverse the session list
- **Launch Session**: Press `Enter` to start the selected session
- **Rename Session**: Press `r` or `R` to rename the selected session
- **Delete Session**: Press `d` to delete the selected session
- **Delete All Sessions**: Press `D` to delete all sessions present.
- **Quit Application**: Press `q` or `Q` to quit the application
- **Save Session**: use the following command to save your session: `kitty-sesh ss <nameofsession>` . if you do not provide the name, a name will be auto generated for you.

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE. See the LICENSE file for details.

## Troubleshooting

The project is still under works so there can be occasional bugs, please report them in the issues, make sure to go through the existing open issues before creating a new one.

## Uninstallation

to uninstall simply delete the kitty-sesh file present at `/usr/local/bin`

```sh
sudo rm /usr/local/bin/kitty-sesh
```
## TODO

- Work on Layouts Resizing
  - ~~Horizontal~~
  - ~~Vertical~~
  - Grid
  - Split
  - Tall
  - Fat
  - ~~Stack~~
- Replace Flex with Grid
- ~~Add the option to delete all sessions using 'D'~~
- ~~Add confirm Modal for delete operation~~
