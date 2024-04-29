#!/usr/bin/env expect

set timeout 300

# Define the package name as a parameter

# Define the command to install the package
set install_command "pacman -Su"

# Spawn a shell
spawn /bin/sh

# Expect the command prompt
expect "# "

# Send the install command
send "$install_command\r"

# Expect the confirmation prompt

expect {
    "Starting full system upgrade" {
        expect {
            "Proceed with installation?" {
                send "y\r"

                expect {
                    "Running post-transaction hooks" {
                        puts "System upgrade complete!"
                    }
                    timeout {
                        puts "Timeout occurred while waiting for post-transaction hooks"
                    }
                }

                exp_continue
            }

            "there is nothing to do" {
                puts "No action taken"
            }

            "error:" {
                puts "Error upgrading system: $expect_out(buffer)"
            }

            timeout {
                puts "Timeout occurred while upgrading system"
            }
        }
    }
}

expect "#"
# Exit the script
send "exit\r"
expect eof