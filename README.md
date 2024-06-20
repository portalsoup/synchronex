# Synchronex

## Overview

This document describes the schema for configuring a user's home environment using HashiCorp Configuration Language (HCL) 
for Synchronex. The schema includes definitions for user-specific settings, package installations, and file operations. 
Synchronex configuration files use the extension .nex.hcl.

## Schema

Synchronex files use the extension .nex.hcl

### Context Configuration

```hcl
context {
  user         = "<username>"
  require_root = <boolean>
}
```

- `user`: The username for which the configuration is being set up.  Used to establish file ownership.  This should be a string.
- `require_root`: A boolean flag indicating if root privileges are required for the setup.

### Provisioner Configuration

The `provisioner` block allows defining one or more actions to set up the user's environment. This includes installing packages and synchronizing configuration files.

```hcl
provisioner "<provisioner_name>" {
    ...
}
```

#### Package Configuration

Within the `provisioner`, each `package` block specifies a software package that the target system requires. On each run, supported
package managers will be automatically detected on the system.  Currently only software managed by a package manager are supported
by this tool.  The attributes for each package are as follows:

- `pacman`: A boolean indicating if the package can be found using `pacman` if present.
- `dpkg`: A boolean indicating if the package can be found using `dpkg` if present.
- `constraints`: A version [range constraint](https://maven.apache.org/enforcer/enforcer-rules/versionRanges.html) for the package, specified as a string.

Example:
```hcl
package "<package_name>" {
    pacman = <boolean>
    dpkg = <boolean>
    constraints = "<version_constraint>"
}
```

#### File Configuration

Each `file` block within the `provisioner` specifies a file operation to be performed. The type of operation can be `sync`, `put`, or `remove`.

- `sync`: Checks for equality of the source and destination files and replaces the destination file if they differ.
- `put`: A one-time initialization, copying the source file to the destination only if the destination is not already present.
- `remove`: Deletes an existing unmanaged file at the destination path.

Attributes for each file operation are as follows:

- `src`: Source file path as a string (required for `sync` and `put` operations).
- `pre_command` (optional): Command to be executed before the file operation, specified as a string.
- `post_command` (optional): Command to be executed after the file operation, specified as a string.

- `owner` (optional): Override the owner of the destination file, specified as a string.
- `group` (optional): Override the group of the destination file, specified as a string.

Additionally, file properties support basic substitution of the global `user` variable using the `{{USER}}` token

Example:
```hcl
file sync "<destination_path>" {
    src = "<source_path>"
    post_command = "<command>"
}
```

### Usage

This schema is intended to be used with Synchronex, which supports HCL configuration files with the .nex.hcl extension. 
It sets up a user's home environment by performing operations based on the provided schema.

Synchronex will scan the working directory for all .nex.hcl files and run them in an undetermined order.  One or more
.nex.hcl files can be passed as program arguments to disable scanning and run one or more specific provisioner

### Example

To use this schema, create an HCL configuration file (e.g., `setup_zsh.nex.hcl`) with the desired settings and run it.

```sh
make
synchronex setup_zsh.nex.hcl
```

```hcl
context {
  user = "myuser"
}

provisioner "setup_zsh" {

  # Require presence of zsh
  package "zsh" {
    dpkg = true
    constraints = "[5, )" # require zsh major version 5 or greater
  }

  # Copy the .zshrc file to the home folder
  file sync "/home/{{USER}}/.zshrc" {
    src = "resources/.zshrc"
    shell = "zsh" # use zsh to run pre/post scripts instead of default bash
    post_command = "source /home/{{USER}}/.zshrc"
  }
}

```

Ensure you have the required permissions and dependencies installed before running the configuration.
