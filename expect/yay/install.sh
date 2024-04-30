#!/usr/bin/env expect

set timeout 300

# Define the package name as a parameter
set package_name [lindex $argv 0]

# Define the command to install the package
set install_command "yay -S $package_name"

# Spawn a shell
spawn /bin/bash

# Expect the command prompt
expect "$ "

# Send the install command
send "$install_command\r"

# Expect the confirmation prompt
expect {
    "Proceed with installation?" {
        # Send "y" to confirm
        send "y\r"
        exp_continue
    }
    "Running post-transaction hooks" {
        # Installation completed successfully, cleaning up
        puts "Package $package_name installed successfully!"
    }

    "$ " {
        # Command exited, back to shell
        exp_continue
    }

    "error:" {
        # There was an error during installation
        puts "Error installing $package_name: $expect_out(buffer)"
    }
    timeout {
        # Timeout occurred
        puts "Timeout occurred while installing $package_name"
    }
}

# Exit the script
send "exit\r"
expect eof